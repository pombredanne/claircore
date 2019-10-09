package suse

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/quay/claircore"
	"github.com/quay/claircore/libvuln/driver"
	"github.com/quay/claircore/pkg/ovalutil"

	"github.com/quay/goval-parser/oval"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Release indicates the SUSE release OVAL database to pull from.
type Release string

// These are some known Releases.
const (
	EnterpriseServer15  = `suse.linux.enterprise.server.15`
	EnterpriseDesktop15 = `suse.linux.enterprise.desktop.15`
	Enterprise15        = `suse.linux.enterprise.15`
	EnterpriseServer12  = `suse.linux.enterprise.server.12`
	EnterpriseDesktop12 = `suse.linux.enterprise.desktop.12`
	Enterprise12        = `suse.linux.enterprise.12`
	EnterpriseServer11  = `suse.linux.enterprise.server.11`
	EnterpriseDesktop11 = `suse.linux.enterprise.desktop.11`
	OpenStackCloud9     = `suse.openstack.cloud.9`
	OpenStackCloud8     = `suse.openstack.cloud.8`
	OpenStackCloud7     = `suse.openstack.cloud.7`
	Leap151             = `opensuse.leap.15.1`
	Leap150             = `opensuse.leap.15.0`
	Leap423             = `opensuse.leap.42.3`
)

var upstreamBase *url.URL

func init() {
	const base = `http://ftp.suse.com/pub/projects/security/oval/`
	var err error
	upstreamBase, err = url.Parse(base)
	if err != nil {
		panic("static url somehow didn't parse")
	}
}

// Updater implements driver.Updater for SUSE.
type Updater struct {
	release string
	ovalutil.Fetcher
	logger *zerolog.Logger
}

var (
	_ driver.Updater   = (*Updater)(nil)
	_ driver.Fetcher   = (*Updater)(nil)
	_ driver.FetcherNG = (*Updater)(nil)
)

// NewUpdater configures an updater to fetch the specified Release.
func NewUpdater(r Release, opts ...Option) (*Updater, error) {
	u := &Updater{
		release: string(r),
	}
	for _, o := range opts {
		if err := o(u); err != nil {
			return nil, err
		}
	}
	if u.logger == nil {
		u.logger = &log.Logger
	}
	l := u.logger.With().Str("component", u.Name()).Logger()
	u.logger = &l
	if u.Fetcher.Client == nil {
		u.Fetcher.Client = http.DefaultClient
	}
	if u.Fetcher.URL == nil {
		var err error
		u.Fetcher.URL, err = upstreamBase.Parse(u.release + ".xml")
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

// Option configures an Updater.
type Option func(*Updater) error

// WithURL overrides the default URL to fetch an OVAL database.
func WithURL(uri, compression string) Option {
	c, cerr := ovalutil.ParseCompressor(compression)
	u, uerr := url.Parse(uri)
	return func(up *Updater) error {
		// Return any errors from the outer function.
		switch {
		case cerr != nil:
			return cerr
		case uerr != nil:
			return uerr
		}
		up.Fetcher.Compression = c
		up.Fetcher.URL = u
		return nil
	}
}

// WithClient sets an http.Client for use with an Updater.
//
// If this Option is not supplied, http.DefaultClient will be used.
func WithClient(c *http.Client) Option {
	return func(u *Updater) error {
		u.Fetcher.Client = c
		return nil
	}
}

// WithLogger sets the default logger.
//
// Functions that take a context.Context will use the logger embedded in there
// instead of the Logger passed in via this Option.
func WithLogger(l *zerolog.Logger) Option {
	return func(u *Updater) error {
		u.logger = l
		return nil
	}
}

// Name satisfies driver.Updater.
func (u *Updater) Name() string {
	return fmt.Sprintf(`suse-updater-%s`, u.release)
}

// Fetch satisfies driver.Fetcher.
func (u *Updater) Fetch() (io.ReadCloser, string, error) {
	ctx := u.logger.WithContext(context.Background())
	ctx, done := context.WithTimeout(ctx, 8*time.Minute)
	defer done()
	r, f, err := u.Fetcher.FetchContext(ctx, driver.Fingerprint(""))
	return r, string(f), err
}

// Parse satisifies the driver.Updater interface.
func (u *Updater) Parse(r io.ReadCloser) ([]*claircore.Vulnerability, error) {
	ctx := u.logger.WithContext(context.Background())
	ctx, done := context.WithTimeout(ctx, 5*time.Minute)
	defer done()
	return u.ParseContext(ctx, r)
}

// ParseContext is like Parse, but with context.
func (u *Updater) ParseContext(ctx context.Context, r io.ReadCloser) ([]*claircore.Vulnerability, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("starting parse")
	defer r.Close()
	root := oval.Root{}
	if err := xml.NewDecoder(r).Decode(&root); err != nil {
		return nil, fmt.Errorf("suse: unable to decode OVAL document: %w", err)
	}
	return ovalutil.NewRPMInfo(&root).Extract(ctx)
}

package dpkg

import (
	"bytes"
	"fmt"

	"github.com/quay/claircore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tadasv/go-dpkg"
)

const (
	osReleasePath = "etc/os-release"
	statusPath    = "var/lib/dpkg/status"
	name          = "dpkg"
	kind          = "package"
	version       = "v0.0.1"
)

// the dpkg scanner is responsible inventorying each package found in the /var/lib/dpkg/status file.
// this scanner should cross reference the way particular vulnerability databases
// classify their CVE's.

// PackageScanner implements the libscan/internal/scanner.PackageScanner interface
// this scanner searches the /var/lib/dpkg/status file for package information.
type PackageScanner struct {
	// the layer hash we are currently scanning
	hash string
	// the status file located in the layer.
	status []byte
	// a logger with context
	logger zerolog.Logger
}

func NewPackageScanner() *PackageScanner {
	ps := &PackageScanner{}
	return ps
}

func (ps *PackageScanner) Name() string {
	return name
}

func (ps *PackageScanner) Version() string {
	return version
}

func (ps *PackageScanner) Kind() string {
	return kind
}

func (ps *PackageScanner) Scan(layer *claircore.Layer) ([]*claircore.Package, error) {
	// set logger context
	ps.logger = log.With().Str("component", "package_scanner").Str("name", ps.Name()).Str("version", ps.Version()).Str("kind", ps.Kind()).Str("layer", layer.Hash).Logger()

	ps.logger.Debug().Msgf("starting scan of layer %v", layer.Hash)

	// scanner maybe shared between layers. reset fields
	ps.hash = layer.Hash

	// extract os-release and dpkg status file
	files, err := layer.Files([]string{osReleasePath, statusPath})
	if err != nil {
		ps.logger.Error().Msgf("searching for files within layer failed: %v", err)
		return nil, fmt.Errorf("searching for files within layer failed: %v", err)
	}

	// add file []byte to PackageScanner
	ps.status = files[statusPath]

	pkgs, err := ps.parsePackages()
	if err != nil {
		ps.logger.Error().Msgf("%v", err)
		return []*claircore.Package{}, err
	}

	log.Printf("dpkg-scanner: done scanning layer %v", layer.Hash)
	return pkgs, nil
}

func (ps *PackageScanner) parsePackages() ([]*claircore.Package, error) {
	// if dpkg status file not found return 0 packages
	if len(ps.status) == 0 {
		ps.logger.Info().Msg("layer did not contain a dpkg status file")
		return []*claircore.Package{}, nil
	}

	reader := bytes.NewReader(ps.status)

	// create dpkg parser
	parser := dpkg.NewParser(reader)
	parsedPkgs := parser.Parse()
	ccPkgs := []*claircore.Package{}

	for _, dpkgPkg := range parsedPkgs {
		ccPkg := &claircore.Package{
			Name:    dpkgPkg.Package,
			Version: dpkgPkg.Version,
			Kind:    "binary",
		}

		if dpkgPkg.Source != "" {
			ccPkg.Source = &claircore.Package{
				Name: dpkgPkg.Source,
				Kind: "source",
				// right now this is an assumption that discovered source package
				// packages relate to their binary versions. we see this in debian
				Version: dpkgPkg.Version,
			}
		}

		ccPkgs = append(ccPkgs, ccPkg)
	}

	return ccPkgs, nil
}

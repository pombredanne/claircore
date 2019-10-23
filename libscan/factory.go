package libscan

import (
	"fmt"

	"github.com/quay/claircore/internal/scanner"
	"github.com/quay/claircore/internal/scanner/defaultfetcher"
	"github.com/quay/claircore/internal/scanner/defaultlayerscanner"
	"github.com/quay/claircore/internal/scanner/defaultscanner"
	"github.com/quay/claircore/pkg/distlock"
	dlpg "github.com/quay/claircore/pkg/distlock/postgres"
)

// ScannerFactory is a factory method to return a Scanner interface during libscan runtime.
type ScannerFactory func(lib *libscan, opts *Opts) (scanner.Scanner, error)

// scannerFactory is the default ScannerFactory
func scannerFactory(lib *libscan, opts *Opts) (scanner.Scanner, error) {
	// add other distributed locking implementations here as they grow
	var sc distlock.Locker
	switch opts.ScanLock {
	case PostgresSL:
		sc = dlpg.NewLock(lib.db, opts.ScanLockRetry)
	default:
		return nil, fmt.Errorf("provided ScanLock opt is unsupported")
	}

	// add other fetcher implementations here as they grow
	var ft scanner.Fetcher
	ft = defaultfetcher.New(lib.client, nil, opts.LayerFetchOpt)

	// convert libscan.Opts to scanner.Opts
	sOpts := &scanner.Opts{
		Store:      lib.store,
		ScanLock:   sc,
		Fetcher:    ft,
		Ecosystems: opts.Ecosystems,
		Vscnrs:     lib.vscnrs,
	}

	// add other layer scanner implementations as they grow
	sOpts.LayerScanner = defaultlayerscanner.New(opts.LayerScanConcurrency, sOpts)
	s := defaultscanner.New(sOpts)
	return s, nil
}

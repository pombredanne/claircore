package dpkg

import (
	"context"

	"github.com/quay/claircore/internal/scanner"
	"github.com/quay/claircore/osrelease"
)

// NewEcosystem provides the set of scanners and coalescers for the dpkg ecosystem
func NewEcosystem(ctx context.Context) *scanner.Ecosystem {
	return &scanner.Ecosystem{
		PackageScanners: []scanner.PackageScannerFunc{
			func(ctx context.Context) (scanner.PackageScanner, error) { return &Scanner{}, nil },
		},
		DistributionScanners: []scanner.DistributionScannerFunc{
			func(ctx context.Context) (scanner.DistributionScanner, error) { return &osrelease.Scanner{}, nil },
		},
		Coalescer: func(ctx context.Context, store scanner.Store) (scanner.Coalescer, error) {
			return NewCoalescer(store), nil
		},
	}
}

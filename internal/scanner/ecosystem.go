package scanner

import "context"

// PackageScannerFunc is a factory method for a scanner.PackageScanner
type PackageScannerFunc func(ctx context.Context) (PackageScanner, error)

// DistributionScannerFunc is a factory method for a scanner.DistributionScanner
type DistributionScannerFunc func(ctx context.Context) (DistributionScanner, error)

// RepositoryScannerFunc is a factory method for a scanner.RepositoryScanner
type RepositoryScannerFunc func(ctx context.Context) (RepositoryScanner, error)

// CoalescerFunc is a factory method for a scanner.Coalescer
type CoalescerFunc func(ctx context.Context, store Store) (Coalescer, error)

// Ecosystems group scanners and coalescers used in a common ecosystem.
//
// Examples of common ecosystems are "dpkg", "rpm", "apk"
type Ecosystem struct {
	PackageScanners      []PackageScannerFunc
	DistributionScanners []DistributionScannerFunc
	RepositoryScanners   []RepositoryScannerFunc
	Coalescer            CoalescerFunc
}

package scanner

import "github.com/quay/claircore/pkg/distlock"

// Opts are options to instantiate a scanner
type Opts struct {
	Store                Store
	ScanLock             distlock.Locker
	LayerScanner         LayerScanner
	Fetcher              Fetcher
	PackageScanners      []PackageScanner
	DistributionScanners []DistributionScanner
	RepositoryScanners   []RepositoryScanner
	Vscnrs               VersionedScanners
}

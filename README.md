# ClairCore

ClairCore provides a set of go modules which handle scanning container layers for installed packages and reporting any discovered vulnerabilities.
ClairCore is designed to be embedded into a service wrapper.

# Usage
Two go modules exist `libscan` and `libvuln`.
These modules export the methods for scanning and image for packages and matching the results of the scan to vulnerabilities respectively. 

## libscan usage
The libscan module exports a single interface
```
type Libscan interface {
	// Scan performs an async scan of a manifest and produces a claircore.ScanReport.
	// Errors encountered before scan begins are returned in the error variable.
	// Errors encountered during scan are populated in the Err field of the claircore.ScanReport
	Scan(ctx context.Context, manifest *claircore.Manifest) (ResultChannel <-chan *claircore.ScanReport, err error)
	// ScanReport tries to retrieve a claircore.ScanReport given the image hash.
	// bool informs caller if found.
	ScanReport(hash string) (*claircore.ScanReport, bool, error)
}
```
Creating an instance
```
opts := &libscan.Opts{
    DataStore: libscan.Postgres,
    ConnString: "postgres://host:port",
    ScanLock: libscan.PostgresSL,
    // see definition for more configuration options
}
lib := libscan.New(opts)
``` 

## libvuln usage
The libvuln module exports a single interface
```
type Libvuln interface {
	Scan(ctx context.Context, sr *claircore.ScanReport) (*claircore.VulnerabilityReport, error)

```
creating an instance
```
opts := &libvuln.Opts{
    DataStore: libvuln.Postgres,
    ConnString: "postgres://host:port",
    ScanLock: libscan.PostgresSL,
    // see definition for more configuration option
}
lib := libvuln.New(opts)
```
Libvuln will first initialize all updaters before returning from its constructor.
Controlling how many updaters initialize in parallel is provided via the libvuln.Opts struct

# Local development and testing
The included makefile has targets for spawning a libscan and a libvuln database instance.
```
make libscan-db-restart
make libvuln-db-restart
```

After both targets are ran you may run integration tests
```
make integration
```

unit tests do not require a database
```
make unit
```

## running libscan on darwin/MacOS
Layer stacking will fail on Darwin/MacOS with a file permissions issue and subsequently fail the scanner. 
In order to get around this the layer stacking integration test has a build tag "unix".
The makefile target `integration-unix` runs tests that will only pass on a unix env. 
You may use the make target 'docker-shell' to drop into a linux shell where `make integration-unix` may be ran.

# Deeper dives
    [Vulnerability Matching](./docs/matching_vulns.md)
    [Vulnerability Tombstoning](./docs/tombstoning.md)


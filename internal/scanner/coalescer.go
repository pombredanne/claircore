package scanner

import (
	"context"

	"github.com/quay/claircore"
)

// layerArifact aggregates the any artifacts found within a layer
type layerArtifacts struct {
	hash  string
	pkgs  []*claircore.Package
	dist  []*claircore.Distribution // each layer can only have a single distribution
	repos []*claircore.Repository
}

// Coalescer takes individual layer artifacts and coalesces them to form the final image's
// package results
//
// It is expected to run a coalescer per "ecosystem". For example it would make sense to coalesce results
// for dpkg, os-release, and apt scanners
type Coalescer struct {
	// a store to access scanartifacts
	store Store
	ps    PackageScanner
	ds    DistributionScanner
	rs    RepositoryScanner
	SR    *claircore.ScanReport
}

// NewCoalescer is a constructor for a Coalescer
func NewCoalescer(store Store, ps PackageScanner, ds DistributionScanner, rs RepositoryScanner) *Coalescer {
	return &Coalescer{
		store: store,
		ps:    ps,
		ds:    ds,
		rs:    rs,
		SR: &claircore.ScanReport{
			PackageIntroduced:     map[int]string{},
			Packages:              map[int]*claircore.Package{},
			Distributions:         map[int]*claircore.Distribution{},
			Repositories:          map[int]*claircore.Repository{},
			DistributionByPackage: map[int]int{},
			RepositoryByPackage:   map[int]int{},
		},
	}
}

// Do coalesces artifacts found in layers and creates a final ScanReport with
// the final package details found in the image. This method blocks and when its finished
// the c.SR field will hold the final ScanReport
func (c *Coalescer) Do(ctx context.Context, layers []*claircore.Layer) error {
	var err error
	// populate layer artifacts
	artifacts := []layerArtifacts{}
	for _, layer := range layers {
		a := layerArtifacts{
			hash: layer.Hash,
		}

		a.pkgs, err = c.store.PackagesByLayer(ctx, layer.Hash, VersionedScanners{c.ps})
		if err != nil {
			return err
		}

		a.dist, err = c.store.DistributionsByLayer(ctx, layer.Hash, VersionedScanners{c.ds})
		if err != nil {
			return err
		}

		a.repos, err = c.store.RepositoriesByLayer(ctx, layer.Hash, VersionedScanners{c.rs})
		if err != nil {
			return err
		}
		artifacts = append(artifacts, a)
	}
	c.associate(ctx, artifacts)
	c.prune(ctx, artifacts)
	return ctx.Err()
}

// Associate searches layer artifacts and records the layer a package was introduced in
// along with the distribution a package is associated with.
func (c *Coalescer) associate(ctx context.Context, artifacts []layerArtifacts) {
	if ctx.Err() != nil {
		return
	}

	var currDist claircore.Distribution

	// search for initial distribution. we will assume that if we find a
	// distribution in layer n all packages in layers 0-n are associated with
	// that distribution
	for _, a := range artifacts {
		if len(a.dist) != 0 {
			currDist = *a.dist[0]
			c.SR.Distributions[currDist.ID] = &currDist
			break
		}
	}

	// associate all packages and handle finding subsequent distributions
	// a subsequent distribution may occur with dist upgrade or downgrade
	for _, a := range artifacts {
		if len(a.dist) != 0 {
			currDist = *a.dist[0]
			c.SR.Distributions[currDist.ID] = &currDist
		}
		for _, pkg := range a.pkgs {
			if _, ok := c.SR.PackageIntroduced[pkg.ID]; !ok {
				c.SR.PackageIntroduced[pkg.ID] = a.hash
			}
			c.SR.Packages[pkg.ID] = pkg
			c.SR.DistributionByPackage[pkg.ID] = currDist.ID
		}
	}

}

// prune removes packages which do not exist in the "newest" verion of a package database.
//
// this assumes we are working with linux distribution package managers and an addition or
// removal of a package presents the entire database in a layer due to CoW semantics.
func (c *Coalescer) prune(ctx context.Context, artifacts []layerArtifacts) {
	if ctx.Err() != nil {
		return
	}

	seenDB := map[string]struct{}{}
	keep := map[string]map[int]struct{}{}

	// walk layer artifacts backwards searching for newest package databases
	for i := len(artifacts) - 1; i >= 0; i-- {
		a := artifacts[i]
		if len(a.pkgs) == 0 {
			continue
		}

		// split package array into a per-package-db representation
		packageDBs := map[string][]*claircore.Package{}
		for _, pkg := range a.pkgs {
			packageDBs[pkg.PackageDB] = append(packageDBs[pkg.PackageDB], pkg)
		}

		// for each package db discovered check if we've seen it. if we haven't
		// seen it record packages in the keep map and add to seenDB
		for db, packages := range packageDBs {
			if _, ok := seenDB[db]; !ok {
				for _, pkg := range packages {
					keep[db][pkg.ID] = struct{}{}
				}
				seenDB[db] = struct{}{}
			}
		}
	}

	// prune any packages not in the keep map
	for id, pkg := range c.SR.Packages {
		if _, ok := keep[pkg.PackageDB][id]; !ok {
			delete(c.SR.Packages, id)
			delete(c.SR.DistributionByPackage, id)
			delete(c.SR.RepositoryByPackage, id)
			delete(c.SR.PackageIntroduced, id)
		}
	}
}

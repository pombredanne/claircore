package controller

import (
	"context"

	"github.com/quay/claircore"
	"github.com/quay/claircore/internal/scanner"
)

// reduce determines which layers should be fetched/scanned and returns these layers
func reduce(ctx context.Context, store scanner.Store, scnrs scanner.VersionedScanners, layers []*claircore.Layer) ([]*claircore.Layer, error) {
	do := []*claircore.Layer{}
	for _, scnr := range scnrs {
		for _, l := range layers {
			if ok, err := store.LayerScanned(ctx, l.Hash, scnr); ok {
				if err != nil {
					return nil, err
				}
				do = append(do, l)
			}
		}
	}
	return do, nil
}

package defaultscanner

import (
	"context"
	"fmt"

	"github.com/quay/claircore"
	"github.com/quay/claircore/internal/scanner"
)

func fetchLayers(ctx context.Context, s *defaultScanner) (ScannerState, error) {
	toFetch := reduce(ctx, s.Store, s.Vscnrs, s.manifest.Layers)
	err := s.Fetcher.Fetch(ctx, toFetch)
	if err != nil {
		s.logger.Error().Str("state", s.getState().String()).Msgf("faild to fetch layers: %v", err)
		return Terminal, fmt.Errorf("failed to fetch layers %v", err)
	}
	return LayerScan, nil
}

// reduce determines which layers should be fetched and returns these layers
func reduce(ctx context.Context, store scanner.Store, scnrs scanner.VersionedScanners, layers []*claircore.Layer) []*claircore.Layer {
	toFetch := []*claircore.Layer{}
	for _, scnr := range scnrs {
		for _, l := range layers {
			if ok, _ := store.LayerScanned(ctx, l.Hash, scnr); ok {
				continue
			}
			toFetch = append(toFetch, l)
		}
	}
	return toFetch
}

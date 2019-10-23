package defaultscanner

import (
	"context"
	"fmt"

	"github.com/quay/claircore"
)

// layerScan is a stateFunc which scans each individual layer using a layerScanner.
// it returns the buildLayerResult stateFunc on success
func layerScan(ctx context.Context, s *defaultScanner) (ScannerState, error) {
	// scan the individual container images. packages will be indexed associated
	// with the individual layer's hash
	s.logger.Debug().Str("state", s.getState().String()).Msg("scanning individual layers in manifest")
	err := s.LayerScanner.Scan(ctx, s.manifest.Hash, s.manifest.Layers)
	if err != nil {
		s.logger.Error().Str("state", s.getState().String()).Msgf("failed to scan all layer contents: %v", err)
		return Terminal, fmt.Errorf("failed to scan all layer contents: %v", err)
	}

	// scan the stacked image layer. packages will be indexed associated with the
	// image's hash.
	s.logger.Debug().Str("state", s.getState().String()).Msg("scanning image layer in manifest")
	err = s.LayerScanner.Scan(ctx, s.manifest.Hash, []*claircore.Layer{s.imageLayer})
	if err != nil {
		s.logger.Error().Str("state", s.getState().String()).Msgf("failed to scan image layer: %v", err)
		return Terminal, fmt.Errorf("failed to scan image layer: %v", err)
	}

	s.logger.Info().Str("state", s.getState().String()).Msg("done scanning layers")
	return Coalesce, nil
}

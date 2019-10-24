package controller

import (
	"context"
	"fmt"
)

func fetchLayers(ctx context.Context, s *Controller) (State, error) {
	toFetch, err := reduce(ctx, s.Store, s.Vscnrs, s.manifest.Layers)
	if err != nil {
		return Terminal, fmt.Errorf("failed to determine layers to fetch: %v", err)
	}
	err = s.Fetcher.Fetch(ctx, toFetch)
	if err != nil {
		s.logger.Error().Str("state", s.getState().String()).Msgf("faild to fetch layers: %v", err)
		return Terminal, fmt.Errorf("failed to fetch layers %v", err)
	}
	return ScanLayers, nil
}

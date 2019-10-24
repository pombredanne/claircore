package controller

import (
	"context"
	"fmt"
)

// scanLayers will run all scanner types against all layers if deemed necssary to scan
func scanLayers(ctx context.Context, c *Controller) (State, error) {
	err := c.LayerScanner.Scan(ctx, c.manifest.Hash, c.manifest.Layers)
	if err != nil {
		return Terminal, fmt.Errorf("failed to scan all layer contents: %v", err)
	}
	return Coalesce, nil
}

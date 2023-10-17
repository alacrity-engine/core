package engine

import (
	"fmt"

	"github.com/alacrity-engine/core/render"
)

var (
	layout  *render.Layout
	batches map[string]*render.Batch
)

func AddBatchToCanvas(canvasName, batchName string, batch *render.Batch, z1, z2 float32) error {
	canvas, err := layout.CanvasByName(canvasName)

	if err != nil {
		return err
	}

	if _, ok := batches[batchName]; ok {
		return fmt.Errorf("batch '%s' already exists", batchName)
	}

	batches[batchName] = batch
	err = canvas.AddBatch(batch, z1, z2)

	if err != nil {
		return err
	}

	return nil
}

func BatchByName(name string) (*render.Batch, error) {
	if batch, ok := batches[name]; ok {
		return nil, fmt.Errorf("batch '%s' doesn't exist", name)
	} else {
		return batch, nil
	}
}

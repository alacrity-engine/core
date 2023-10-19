package engine

import (
	"github.com/alacrity-engine/core/render"
)

var (
	layout *render.Layout
)

func AddBatchToCanvas(canvasName string, batch *render.Batch, z1, z2 float32) error {
	canvas, err := layout.CanvasByName(canvasName)

	if err != nil {
		return err
	}

	if _, err := canvas.BatchByName(batch.Name()); err != nil {
		return err
	}

	err = canvas.AddBatch(batch, z1, z2)

	if err != nil {
		return err
	}

	return nil
}

func BatchByName(canvasID, name string) (*render.Batch, error) {
	canvas, err := layout.CanvasByName(canvasID)

	if err != nil {
		return nil, err
	}

	if batch, err := canvas.BatchByName(name); err != nil {
		return nil, err
	} else {
		return batch, nil
	}
}

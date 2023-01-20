package render

import (
	"fmt"
	"sort"
)

// TODO: add a remove canvas method.

// TODO: change the global variables
// for views and projections on all
// the batches when a new canvas gets added.

type Layout struct {
	zMin     float32
	zMax     float32
	canvases []*Canvas
	batches  []*Batch
}

func (layot *Layout) Range() (float32, float32) {
	return layot.zMin, layot.zMax
}

func (layout *Layout) AddCanvas(canvas *Canvas) error {
	if len(layout.canvases) >= 256 {
		return fmt.Errorf("max number of canvases exceeded")
	}

	if canvas == nil {
		return fmt.Errorf("the canvas is nil")
	}

	length := len(layout.canvases)
	ind := sort.Search(length, func(i int) bool {
		return layout.canvases[i].index >= canvas.index
	})

	if ind < length && layout.canvases[ind].index == canvas.index {
		return fmt.Errorf(
			"the canvas with index %d already exists on the layot", canvas.index)
	}

	if ind == 0 {
		layout.canvases = append(layout.canvases, nil)
		copy(layout.canvases[1:], layout.canvases)
		layout.canvases[0] = canvas

		canvasZMin, canvasZMax := canvas.Range()
		layout.zMin = canvasZMin

		if length == 0 {
			layout.zMax = canvasZMax
		}
	} else if ind < length {
		layout.canvases = append(layout.canvases[:ind+1],
			layout.canvases[ind:]...)
		layout.canvases[ind] = canvas
	} else {
		layout.canvases = append(layout.canvases, canvas)

		canvasZMin, canvasZMax := canvas.Range()
		layout.zMax = canvasZMax

		if length == 0 {
			layout.zMin = canvasZMin
		}
	}

	canvas.layout = layout
	canvas.pos = byte(ind)

	for i := 0; i < len(layout.batches); i++ {
		batch := layout.batches[i]

		batch.setCanvasProjection(canvas.index, canvas.projection)
		batch.setCanvasView(canvas.index, canvas.camera.View())
	}

	return nil
}

func NewLayout() *Layout {
	return &Layout{
		zMin:     0,
		zMax:     0,
		canvases: []*Canvas{},
	}
}

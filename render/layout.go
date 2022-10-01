package render

import (
	"fmt"
	"sort"
)

type Layout struct {
	zMin     float32
	zMax     float32
	canvases []*Canvas
}

func (layot *Layout) Range() (float32, float32) {
	return layot.zMin, layot.zMax
}

func (layout *Layout) AddCanvas(canvas *Canvas) error {
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
		canvases := make([]*Canvas, 0, length+1)
		canvases = append(canvases, canvas)
		canvases = append(canvases, layout.canvases...)
		layout.canvases = canvases

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

	return nil
}

func NewLayout() *Layout {
	return &Layout{
		zMin:     0,
		zMax:     0,
		canvases: []*Canvas{},
	}
}

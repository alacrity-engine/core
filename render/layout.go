package render

import (
	"fmt"
	"sort"
)

// TODO: add a remove canvas method.

type Layout struct {
	zMin      float32
	zMax      float32
	canvases  []*Canvas
	nameIndex map[string]*Canvas
}

func (layot *Layout) Range() (float32, float32) {
	return layot.zMin, layot.zMax
}

func (layout *Layout) Draw() error {
	for _, canvas := range layout.canvases {
		err := canvas.draw()

		if err != nil {
			return err
		}
	}

	return nil
}

func (layout *Layout) AddCanvas(canvas *Canvas) error {
	if _, ok := layout.nameIndex[canvas.name]; ok {
		return fmt.Errorf(
			"a canvas named '%s' already exists on the layout", canvas.name)
	}

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
	layout.nameIndex[canvas.name] = canvas

	return nil
}

func NewLayout() *Layout {
	return &Layout{
		zMin:      0,
		zMax:      0,
		canvases:  []*Canvas{},
		nameIndex: map[string]*Canvas{},
	}
}

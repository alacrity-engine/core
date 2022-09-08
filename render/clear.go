package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type ClearBit uint32

const (
	ClearBitColor   ClearBit = gl.COLOR_BUFFER_BIT
	ClearBitDepth   ClearBit = gl.DEPTH_BUFFER_BIT
	ClearBitStencil ClearBit = gl.STENCIL_BUFFER_BIT
)

func SetClearColor(_color RGBA) {
	gl.ClearColor(float32(_color.R), float32(_color.G), float32(_color.B), float32(_color.A))
}

func Clear(bit ClearBit) {
	gl.Clear(uint32(bit))
}

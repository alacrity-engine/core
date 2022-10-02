package render

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type ClearBit uint32

const (
	ClearBitColor   ClearBit = gl.COLOR_BUFFER_BIT
	ClearBitDepth   ClearBit = gl.DEPTH_BUFFER_BIT
	ClearBitStencil ClearBit = gl.STENCIL_BUFFER_BIT
)

func SetClearColor(_color RGBA) {
	gl.ClearColor(_color.R, _color.G, _color.B, _color.A)
}

func Clear(bit ClearBit) {
	gl.Clear(uint32(bit))
}

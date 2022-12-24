package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func Initialize(_width, _height int, _zMin, _zMax float32) error {
	if _zMin >= _zMax {
		return fmt.Errorf("max Z must be greater tham min Z")
	}

	zMin = _zMin
	zMax = _zMax

	err := SetWidth(_width)

	if err != nil {
		return err
	}

	err = SetHeight(_height)

	if err != nil {
		return err
	}

	err = gl.Init()

	if err != nil {
		return err
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Enable(gl.DEPTH_TEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	spriteIndices := []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
	}

	gl.GenBuffers(1, &spriteIndexBufferHandler)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, spriteIndexBufferHandler)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(spriteIndices)*4, gl.Ptr(spriteIndices), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return nil
}

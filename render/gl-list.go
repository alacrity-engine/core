package render

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"golang.org/x/exp/constraints"
)

type numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

type glList[T numeric] struct {
	glHandler uint32
	length    int
	capactity int
	drawMode  DrawMode
}

func (list *glList[T]) setData(data []T) {
	var zeroVal T
	dataLength := len(data) * int(unsafe.Sizeof(zeroVal))
	defer func() {
		list.length = dataLength
	}()

	if dataLength > list.capactity && list.glHandler != 0 {
		gl.DeleteBuffers(1, &list.glHandler)
		list.glHandler = 0
	}

	if list.glHandler == 0 {
		var glHandler uint32
		gl.GenBuffers(1, &glHandler)
		gl.BindBuffer(gl.ARRAY_BUFFER, glHandler)
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*
			int(unsafe.Sizeof(zeroVal)), gl.Ptr(data),
			uint32(list.drawMode))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		list.glHandler = glHandler
		list.capactity = dataLength

		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, dataLength, gl.Ptr(data))

	if list.capactity-dataLength > 0 {
		gl.BufferSubData(gl.ARRAY_BUFFER, dataLength,
			list.capactity-dataLength, gl.Ptr(nil))
	}
}

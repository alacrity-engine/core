package render

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"golang.org/x/exp/constraints"
)

type numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

type gpuList[T numeric] struct {
	glHandler uint32
	stride    int
	length    int
	capacity  int
	drawMode  DrawMode
}

func (list *gpuList[T]) getLength() int {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	return list.length / dataSize
}

func (list *gpuList[T]) getCapacity() int {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	return list.capacity / dataSize
}

func (list *gpuList[T]) grow(targetCap int) {
	// TODO: come up with a better algorithm
	// for GPU list growth (based on the stride).

	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))
	//growFactor := 2

	if list.capacity > (1<<10)*dataSize {
		//growFactor = 4
	}

	// Allocate a greater buffer for the GPU list.
	adjCap := list.capacity / list.stride
	newCapacity := adjCap * 2
	newCapacity *= list.stride

	if targetCap > newCapacity {
		newCapacity = targetCap
	}

	if newCapacity%list.stride != 0 {
		_ = newCapacity
	}

	var glHandler uint32
	gl.GenBuffers(1, &glHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, glHandler)
	gl.BufferData(gl.ARRAY_BUFFER, newCapacity, gl.Ptr(nil),
		uint32(list.drawMode))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// Copy all the data from the old
	// buffer to the new buffer.
	gl.BindBuffer(gl.COPY_READ_BUFFER, list.glHandler)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, glHandler)
	gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER,
		0, 0, list.length)
	gl.BindBuffer(gl.COPY_READ_BUFFER, 0)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, 0)

	// Delete the old buffer.
	oldBuffer := list.glHandler
	list.glHandler = glHandler
	list.capacity = newCapacity
	gl.DeleteBuffers(1, &oldBuffer)
}

func (list *gpuList[T]) addElement(elem T) {
	if list.glHandler == 0 {
		list.setData([]T{elem})
		return
	}

	dataSize := int(unsafe.Sizeof(elem))

	if list.length+dataSize > list.capacity {
		list.grow(list.length + dataSize)
	}

	defer func() {
		list.length += dataSize
	}()

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, list.length,
		dataSize, gl.Ptr([]T{elem}))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (list *gpuList[T]) addElements(elems []T) {
	if list.glHandler == 0 {
		list.setData(elems)
		return
	}

	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if list.length+len(elems)*dataSize > list.capacity {
		list.grow(list.length + len(elems)*dataSize)
	}

	defer func() {
		list.length += len(elems) * dataSize
	}()

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, list.length,
		len(elems)*dataSize, gl.Ptr(elems))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (list *gpuList[T]) replaceElement(idx int, elem T) error {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if idx*dataSize > list.length-dataSize {
		return fmt.Errorf(
			"wrong index %d with data length %d",
			idx, list.length)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, idx*dataSize,
		dataSize, gl.Ptr([]T{elem}))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

func (list *gpuList[T]) replaceElements(offset, count int, data []T) error {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if (offset+count-1)*dataSize > list.length-dataSize {
		return fmt.Errorf(
			"wrong offset %d and count %d with data length %d",
			offset, count, list.length)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, offset*dataSize,
		count*dataSize, gl.Ptr(data))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

func (list *gpuList[T]) removeElement(idx int) error {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if idx*dataSize > list.length-dataSize {
		return fmt.Errorf(
			"wrong index %d with data length %d",
			idx, list.length)
	}

	gl.BindBuffer(gl.COPY_READ_BUFFER, list.glHandler)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, list.glHandler)
	gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER,
		(idx+1)*dataSize, idx*dataSize, dataSize)
	gl.BindBuffer(gl.COPY_READ_BUFFER, 0)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, 0)

	originalLength := list.length
	list.length -= dataSize

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, originalLength-dataSize,
		dataSize, gl.Ptr(nil))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

func (list *gpuList[T]) clear() {
	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.ClearBufferData(gl.ARRAY_BUFFER, gl.R8UI,
		gl.RED, gl.BYTE, gl.Ptr([]byte{0}))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (list *gpuList[T]) removeElements(offset, count int) error {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if (offset+count-1)*dataSize > list.length-dataSize {
		return fmt.Errorf(
			"wrong offset %d and count %d with data length %d",
			offset, count, list.length)
	}

	gl.BindBuffer(gl.COPY_READ_BUFFER, list.glHandler)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, list.glHandler)
	gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER,
		(offset+1)*dataSize, offset*dataSize, count*dataSize)
	gl.BindBuffer(gl.COPY_READ_BUFFER, 0)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, 0)

	originalLength := list.length
	list.length -= count * dataSize

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, originalLength-count*dataSize,
		count*dataSize, gl.Ptr(nil))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

func (list *gpuList[T]) setData(data []T) {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))
	dataLength := len(data) * dataSize
	defer func() {
		list.length = dataLength
	}()

	if dataLength > list.capacity && list.glHandler != 0 {
		gl.DeleteBuffers(1, &list.glHandler)
		list.glHandler = 0
	}

	if dataLength <= 0 {
		if list.glHandler != 0 {
			gl.DeleteBuffers(1, &list.glHandler)
			list.glHandler = 0
		}

		return
	}

	if list.glHandler == 0 {
		var glHandler uint32
		gl.GenBuffers(1, &glHandler)
		gl.BindBuffer(gl.ARRAY_BUFFER, glHandler)
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*
			dataSize, gl.Ptr(data),
			uint32(list.drawMode))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		list.glHandler = glHandler
		list.length = dataLength
		list.capacity = dataLength

		// Read.
		//dbg := make([]T, list.length/dataSize)
		//gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
		//gl.GetBufferSubData(gl.ARRAY_BUFFER, 0, list.length, gl.Ptr(dbg))
		//
		//_ = dbg

		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, list.glHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, dataLength, gl.Ptr(data))

	if list.capacity-dataLength > 0 {
		gl.BufferSubData(gl.ARRAY_BUFFER, dataLength,
			list.capacity-dataLength, gl.Ptr(nil))
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (list *gpuList[T]) addDataFromBuffer(buffer uint32, count int) {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if list.length+count*dataSize > list.capacity {
		list.grow(list.length + count*dataSize)
	}

	gl.BindBuffer(gl.COPY_READ_BUFFER, buffer)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, list.glHandler)
	gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER,
		0, list.length, count*dataSize)
	gl.BindBuffer(gl.COPY_READ_BUFFER, 0)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, 0)
}

func (list *gpuList[T]) copyDataToBuffer(buffer uint32, offset, count int) error {
	var zeroVal T
	dataSize := int(unsafe.Sizeof(zeroVal))

	if (offset+count-1)*dataSize > list.length-dataSize {
		return fmt.Errorf(
			"wrong offset %d and count %d with data length %d",
			offset, count, list.length)
	}

	gl.BindBuffer(gl.COPY_READ_BUFFER, list.glHandler)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, buffer)
	gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER,
		offset*dataSize, 0, count*dataSize)
	gl.BindBuffer(gl.COPY_READ_BUFFER, 0)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, 0)

	return nil
}

func newGPUList[T numeric](mode DrawMode, initData []T, stride int) *gpuList[T] {
	list := &gpuList[T]{
		drawMode: mode,
		stride:   stride,
	}

	list.setData(initData)

	return list
}

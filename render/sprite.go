package render

import (
	"fmt"
	"unsafe"

	"github.com/alacrity-engine/core/geometry"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type DrawMode uint32

const (
	DrawModeStatic  DrawMode = gl.STATIC_DRAW
	DrawModeDynamic DrawMode = gl.DYNAMIC_DRAW
	DrawModeStream  DrawMode = gl.STREAM_DRAW
)

type Sprite struct {
	glHandler                         uint32
	glVertexBufferHandler             uint32
	glTextureCoordinatesBufferHandler uint32
	texture                           *Texture
	shaderProgram                     *ShaderProgram
	drawMode                          DrawMode
}

func NewSpriteFromTextureAndProgram(drawMode DrawMode, texture *Texture, shaderProgram *ShaderProgram, targetArea geometry.Rect) (*Sprite, error) {
	if texture == nil || texture.glHandler == 0 {
		return nil, fmt.Errorf("no texture")
	}

	if shaderProgram == nil || shaderProgram.glHandler == 0 {
		return nil, fmt.Errorf("no shader program")
	}

	// TODO: draw sprites with original size
	// using relativity to the screen size.

	textureAspectRatio := float32(targetArea.W() / targetArea.H())
	vertices := []float32{
		textureAspectRatio * -1.0, -1.0, 0.0,
		textureAspectRatio * -0.1, 0.1, 0.0,
		textureAspectRatio * 0.1, 0.1, 0.0,
		textureAspectRatio * 0.1, -0.1, 0.0,
	}
	textureCoordinates := []float32{
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight),
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight),
		float32(targetArea.Max.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight),
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight),
	}

	var handler uint32
	gl.GenVertexArrays(1, &handler)
	gl.BindVertexArray(handler)

	var vertexBufferHandler uint32
	gl.GenBuffers(1, &vertexBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, unsafe.Pointer(uintptr(0)))

	var textureCoordinatesBufferHandler uint32
	gl.GenBuffers(1, &textureCoordinatesBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordinatesBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(textureCoordinates)*4, gl.Ptr(textureCoordinates), uint32(drawMode))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, unsafe.Pointer(uintptr(0)))

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return &Sprite{
		glHandler:                         handler,
		glVertexBufferHandler:             vertexBufferHandler,
		glTextureCoordinatesBufferHandler: textureCoordinatesBufferHandler,
		texture:                           texture,
		shaderProgram:                     shaderProgram,
		drawMode:                          drawMode,
	}, nil
}

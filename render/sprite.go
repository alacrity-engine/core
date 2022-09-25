package render

import (
	"fmt"

	"github.com/alacrity-engine/core/geometry"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type DrawMode uint32

const (
	DrawModeStatic  DrawMode = gl.STATIC_DRAW
	DrawModeDynamic DrawMode = gl.DYNAMIC_DRAW
	DrawModeStream  DrawMode = gl.STREAM_DRAW
)

var (
	spriteIndexBufferHandler uint32
)

type Sprite struct {
	glHandler                         uint32
	glVertexBufferHandler             uint32
	glTextureCoordinatesBufferHandler uint32
	texture                           *Texture
	shaderProgram                     *ShaderProgram
	drawMode                          DrawMode
}

func (sprite *Sprite) SetTargetArea(targetArea geometry.Rect) error {
	textureRect := geometry.R(0, 0,
		float64(sprite.texture.imageWidth),
		float64(sprite.texture.imageHeight))

	if !textureRect.Contains(targetArea.Min) || !textureRect.Contains(targetArea.Max) ||
		!textureRect.Contains(geometry.V(targetArea.Min.X, targetArea.Max.Y)) ||
		!textureRect.Contains(geometry.V(targetArea.Max.X, targetArea.Min.Y)) {
		return fmt.Errorf(
			"rectangle '%v' cannot server as a texture subarea for the sprite", targetArea)
	}

	textureCoordinates := []float32{
		float32(targetArea.Min.X) / float32(sprite.texture.imageWidth), float32(targetArea.Min.Y) / float32(sprite.texture.imageHeight),
		float32(targetArea.Min.X) / float32(sprite.texture.imageWidth), float32(targetArea.Max.Y) / float32(sprite.texture.imageHeight),
		float32(targetArea.Max.X) / float32(sprite.texture.imageWidth), float32(targetArea.Min.Y) / float32(sprite.texture.imageHeight), // extraneous 3
		float32(targetArea.Min.X) / float32(sprite.texture.imageWidth), float32(targetArea.Max.Y) / float32(sprite.texture.imageHeight), // extraneous 1
		float32(targetArea.Max.X) / float32(sprite.texture.imageWidth), float32(targetArea.Max.Y) / float32(sprite.texture.imageHeight),
		float32(targetArea.Max.X) / float32(sprite.texture.imageWidth), float32(targetArea.Min.Y) / float32(sprite.texture.imageHeight),
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glTextureCoordinatesBufferHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(textureCoordinates)*4, gl.Ptr(textureCoordinates))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

func (sprite *Sprite) Draw(model, view, projection mgl32.Mat4) {
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, spriteIndexBufferHandler)
	//defer gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	sprite.shaderProgram.Use()
	defer gl.UseProgram(0)

	gl.BindVertexArray(sprite.glHandler)
	defer gl.BindVertexArray(0)

	sprite.texture.Use()
	defer func() {
		gl.ActiveTexture(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}()

	sprite.shaderProgram.SetMatrix4("model", model)
	sprite.shaderProgram.SetMatrix4("view", view)
	sprite.shaderProgram.SetMatrix4("projection", projection)

	//gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func NewSpriteFromTextureAndProgram(drawMode DrawMode, texture *Texture, shaderProgram *ShaderProgram, targetArea geometry.Rect) (*Sprite, error) {
	if texture == nil || texture.glHandler == 0 {
		return nil, fmt.Errorf("no texture")
	}

	if shaderProgram == nil || shaderProgram.glHandler == 0 {
		return nil, fmt.Errorf("no shader program")
	}

	texToScreenWidth := float32(targetArea.W() / float64(width))
	texToscreenHeight := float32(targetArea.H() / float64(width))
	vertices := []float32{
		texToScreenWidth * -1.0, texToscreenHeight * -1.0, 0.0,
		texToScreenWidth * -1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0, // extraneous 3
		texToScreenWidth * -1.0, texToscreenHeight * 1.0, 0.0, // extraneous 1
		texToScreenWidth * 1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0,
	}
	textureCoordinates := []float32{
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight),
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight),
		float32(targetArea.Max.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight), // extraneous 3
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight), // extraneous 1
		float32(targetArea.Max.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight),
		float32(targetArea.Max.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight),
	}

	var handler uint32
	gl.GenVertexArrays(1, &handler)
	gl.BindVertexArray(handler)

	var vertexBufferHandler uint32
	gl.GenBuffers(1, &vertexBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	var textureCoordinatesBufferHandler uint32
	gl.GenBuffers(1, &textureCoordinatesBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordinatesBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(textureCoordinates)*4, gl.Ptr(textureCoordinates), uint32(drawMode))
	texCoordAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

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

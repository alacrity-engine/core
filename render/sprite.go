package render

import (
	"fmt"

	"github.com/alacrity-engine/core/geometry"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	spriteIndexBufferHandler uint32
)

type Sprite struct {
	glHandler                         uint32
	glVertexBufferHandler             uint32
	glTextureCoordinatesBufferHandler uint32
	glColorMaskBufferHandler          uint32
	texture                           *Texture
	shaderProgram                     *ShaderProgram
	drawMode                          DrawMode
	drawZ                             float32 // drawZ must be in the range of [-1; 1]
	canvas                            *Canvas
	batch                             *Batch
	batchIndex                        int
}

func (sprite *Sprite) createVertexBuffer() {
	var vertexBufferHandler uint32
	gl.GenBuffers(1, &vertexBufferHandler)
	sprite.glVertexBufferHandler = vertexBufferHandler
}

func (sprite *Sprite) deleteVertexBuffer() {
	gl.DeleteBuffers(1, &sprite.glVertexBufferHandler)
	sprite.glVertexBufferHandler = 0
}

func (sprite *Sprite) createTextureCoordinatesBuffer() {
	var texCoordBufferHandler uint32
	gl.GenBuffers(1, &texCoordBufferHandler)
	sprite.glTextureCoordinatesBufferHandler = texCoordBufferHandler
}

func (sprite *Sprite) deleteTextureCoordinatesBuffer() {
	gl.DeleteBuffers(1, &sprite.glTextureCoordinatesBufferHandler)
	sprite.glTextureCoordinatesBufferHandler = 0
}

func (sprite *Sprite) createColorMaskBuffer() {
	var colorMaskBufferHandler uint32
	gl.GenBuffers(1, &colorMaskBufferHandler)
	sprite.glColorMaskBufferHandler = colorMaskBufferHandler
}

func (sprite *Sprite) deleteColorMaskBuffer() {
	gl.DeleteBuffers(1, &sprite.glColorMaskBufferHandler)
	sprite.glColorMaskBufferHandler = 0
}

func (sprite *Sprite) createVertexArray() {
	var handler uint32
	gl.GenVertexArrays(1, &handler)
	sprite.glHandler = handler
}

func (sprite *Sprite) deleteVertexArray() {
	gl.DeleteVertexArrays(1, &sprite.glHandler)
	sprite.glHandler = 0
}

func (sprite *Sprite) assembleVertexArray() {
	gl.BindVertexArray(sprite.glHandler)

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glVertexBufferHandler)
	vertAttrib := uint32(gl.GetAttribLocation(sprite.shaderProgram.glHandler, gl.Str("aPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glTextureCoordinatesBufferHandler)
	texCoordAttrib := uint32(gl.GetAttribLocation(sprite.shaderProgram.glHandler, gl.Str("aTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glColorMaskBufferHandler)
	colorAttrib := uint32(gl.GetAttribLocation(sprite.shaderProgram.glHandler, gl.Str("aColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (sprite *Sprite) SetZ(z float32) {
	sprite.drawZ = mgl32.Clamp(z, -1, 1)
}

func (sprite *Sprite) SetColorMask(colorMask [4]RGBA) {
	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glColorMaskBufferHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(colorMask)*4*4, gl.Ptr(colorMask[:]))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (sprite *Sprite) SetTargetArea(targetArea geometry.Rect) error {
	textureRect := geometry.R(0, 0,
		float64(sprite.texture.imageWidth),
		float64(sprite.texture.imageHeight))

	if !textureRect.Contains(targetArea.Min) || !textureRect.Contains(targetArea.Max) ||
		!textureRect.Contains(geometry.V(targetArea.Min.X, targetArea.Max.Y)) ||
		!textureRect.Contains(geometry.V(targetArea.Max.X, targetArea.Min.Y)) {
		return fmt.Errorf(
			"rectangle '%v' cannot serveS as a texture subarea for the sprite", targetArea)
	}

	textureCoordinates := []float32{
		float32(targetArea.Min.X) / float32(sprite.texture.imageWidth), float32(targetArea.Min.Y) / float32(sprite.texture.imageHeight),
		float32(targetArea.Min.X) / float32(sprite.texture.imageWidth), float32(targetArea.Max.Y) / float32(sprite.texture.imageHeight),
		float32(targetArea.Max.X) / float32(sprite.texture.imageWidth), float32(targetArea.Max.Y) / float32(sprite.texture.imageHeight),
		float32(targetArea.Max.X) / float32(sprite.texture.imageWidth), float32(targetArea.Min.Y) / float32(sprite.texture.imageHeight),
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glTextureCoordinatesBufferHandler)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(textureCoordinates)*4, gl.Ptr(textureCoordinates))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

// TODO: we don't need the space shrinking
// in the range conversion formula. Therefore
// no need for it at all. Instead just convert
// the canvas limits from relative to global
// (with the Range() method) and add the sprite
// Z to the leftmost canvas' global limit in
// order to obtain the sprite's global Z. This
// method is more appropriate because a
// canvas may have a perspective projection
// for 3D, and the global canvas space shrinking
// is not appropriate for it.

// TODO: rename the ortho2D projection to the
// ortho2DStandard. Don't hard-code its values
// in the range conversion formula or the canvas
// code. Each canvas may have a different projection.

func (sprite *Sprite) draw(model, view, projection mgl32.Mat4) {
	if sprite.batch != nil {
		return
	}

	sprite.shaderProgram.Use()
	defer gl.UseProgram(0)

	gl.BindVertexArray(sprite.glHandler)
	defer gl.BindVertexArray(0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, spriteIndexBufferHandler)
	defer gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	sprite.texture.Use()
	defer func() {
		gl.ActiveTexture(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}()

	zModifier := sprite.drawZ

	if sprite.canvas != nil && sprite.canvas.layout != nil {
		globalZMin, globalZMax := sprite.canvas.layout.Range()
		// A range conversion formula.
		zModifier = 2*(sprite.drawZ+sprite.canvas.Z()-globalZMin)/
			(globalZMax-globalZMin) - 1
	}

	model[12] /= float32(width)
	model[13] /= float32(width)
	model[14] += zModifier

	view[12] /= -float32(width)
	view[13] /= -float32(width)

	//debug := projection.Mul4(view).Mul4(model)
	//_ = debug

	sprite.shaderProgram.SetMatrix4("model", model)
	sprite.shaderProgram.SetMatrix4("view", view)
	sprite.shaderProgram.SetMatrix4("projection", projection)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func (sprite *Sprite) drawToBatch(transform *geometry.Transform) {
	// TODO: copy all the dynamic data to
	// the batch buffers.
}

func (sprite *Sprite) Draw(transform *geometry.Transform) error {
	if transform == nil {
		return fmt.Errorf("the transform is nil")
	}

	if sprite.canvas == nil {
		return fmt.Errorf("the sprite has no canvas")
	}

	if sprite.batch != nil {
		sprite.drawToBatch(transform)

		return nil
	}

	sprite.draw(transform.Data(), sprite.canvas.
		camera.View(), Ortho2DStandard())

	return nil
}

func NewSpriteFromTextureAndProgram(textureDrawMode, colorDrawMode DrawMode, texture *Texture, shaderProgram *ShaderProgram, targetArea geometry.Rect) (*Sprite, error) {
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
		texToScreenWidth * 1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0,
	}
	textureCoordinates := []float32{
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight),
		float32(targetArea.Min.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight),
		float32(targetArea.Max.X) / float32(texture.imageWidth), float32(targetArea.Max.Y) / float32(texture.imageHeight),
		float32(targetArea.Max.X) / float32(texture.imageWidth), float32(targetArea.Min.Y) / float32(texture.imageHeight),
	}
	colorMask := []float32{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
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
	gl.BufferData(gl.ARRAY_BUFFER, len(textureCoordinates)*4, gl.Ptr(textureCoordinates), uint32(textureDrawMode))
	texCoordAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	var colorMaskBufferHandler uint32
	gl.GenBuffers(1, &colorMaskBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorMaskBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(colorMask)*4, gl.Ptr(colorMask), uint32(colorDrawMode))
	colorAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return &Sprite{
		glHandler:                         handler,
		glVertexBufferHandler:             vertexBufferHandler,
		glTextureCoordinatesBufferHandler: textureCoordinatesBufferHandler,
		glColorMaskBufferHandler:          colorMaskBufferHandler,
		texture:                           texture,
		shaderProgram:                     shaderProgram,
		drawMode:                          textureDrawMode,
		batchIndex:                        -1,
	}, nil
}

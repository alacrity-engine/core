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
	colorMask                         ColorMask
	targetArea                        geometry.Rect
	texture                           *Texture
	shaderProgram                     *ShaderProgram
	drawMode                          DrawMode
	drawZ                             float32 // drawZ must be in the range of [zMin; zMax]
	canvas                            *Canvas
	batch                             *Batch
	batchIndex                        int
}

func (sprite *Sprite) ColorMask() ColorMask {
	return sprite.colorMask
}

func (sprite *Sprite) TargetArea() geometry.Rect {
	return sprite.targetArea
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
	sprite.drawZ = mgl32.Clamp(z, zMin, zMax)
}

func (sprite *Sprite) SetColorMask(colorMask ColorMask) error {
	data := colorMask.Data()

	if sprite.batch == nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glColorMaskBufferHandler)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(colorMask)*4*4, gl.Ptr(data[:]))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		return nil
	}

	colorMaskData := make([]float32, 24)
	geometry.ColorMaskDataNoElementsFill(colorMaskData, data)
	err := sprite.batch.colorMasks.replaceElements(
		sprite.batchIndex, len(colorMaskData), colorMaskData)

	if err != nil {
		return err
	}

	sprite.colorMask = colorMask

	return nil
}

func (sprite *Sprite) SetTargetArea(targetArea geometry.Rect) error {
	textureRect := geometry.R(0, 0,
		float64(sprite.texture.imageWidth),
		float64(sprite.texture.imageHeight))

	if !textureRect.Contains(targetArea.Min) || !textureRect.Contains(targetArea.Max) ||
		!textureRect.Contains(geometry.V(targetArea.Min.X, targetArea.Max.Y)) ||
		!textureRect.Contains(geometry.V(targetArea.Max.X, targetArea.Min.Y)) {
		return fmt.Errorf(
			"rectangle '%v' cannot serve as a texture subarea for the sprite", targetArea)
	}

	if sprite.batch == nil {
		textureCoordinates := make([]float32, 8)
		geometry.ComputeSpriteTextureCoordinatesFill(textureCoordinates,
			sprite.texture.imageWidth, sprite.texture.imageHeight, targetArea)

		gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glTextureCoordinatesBufferHandler)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(textureCoordinates)*4, gl.Ptr(textureCoordinates))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		vertices := make([]float32, 12)
		geometry.ComputeSpriteVerticesFill(vertices, width, height, targetArea)

		gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glVertexBufferHandler)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		return nil
	}

	textureCoordinates := make([]float32, 12)
	geometry.ComputeSpriteTextureCoordinatesNoElementsFill(textureCoordinates,
		sprite.texture.imageWidth, sprite.texture.imageHeight, targetArea)
	err := sprite.batch.texCoords.replaceElements(sprite.batchIndex,
		len(textureCoordinates), textureCoordinates)

	if err != nil {
		return err
	}

	vertices := make([]float32, 18)
	geometry.ComputeSpriteVerticesNoElementsFill(
		vertices, width, height, targetArea)
	err = sprite.batch.vertices.replaceElements(
		sprite.batchIndex, len(vertices), vertices)

	if err != nil {
		return err
	}

	sprite.targetArea = targetArea

	return nil
}

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
		zModifier = (zMax-zMin)*(sprite.drawZ+sprite.canvas.Z()-globalZMin)/
			(globalZMax-globalZMin) + zMin
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

func (sprite *Sprite) drawToBatch(model mgl32.Mat4) error {
	zModifier := sprite.drawZ

	if sprite.canvas != nil && sprite.canvas.layout != nil {
		globalZMin, globalZMax := sprite.canvas.layout.Range()
		// A range conversion formula.
		zModifier = (zMax-zMin)*(sprite.drawZ+sprite.canvas.Z()-globalZMin)/
			(globalZMax-globalZMin) + zMin
	}

	model[12] /= float32(width)
	model[13] /= float32(width)
	model[14] += zModifier

	err := sprite.batch.models.replaceElements(
		sprite.batchIndex, len(model), model[:])

	if err != nil {
		return err
	}

	err = sprite.batch.shouldDraw.replaceElement(
		sprite.batchIndex, 1)

	if err != nil {
		return err
	}

	return nil
}

func (sprite *Sprite) Draw(transform *geometry.Transform) error {
	if transform == nil {
		return fmt.Errorf("the transform is nil")
	}

	if sprite.canvas == nil {
		return fmt.Errorf("the sprite has no canvas")
	}

	if sprite.batch != nil {
		return sprite.drawToBatch(transform.Data())
	}

	sprite.draw(transform.Data(), sprite.canvas.
		camera.View(), sprite.canvas.projection)

	return nil
}

func NewSpriteFromTextureAndProgram(textureDrawMode, colorDrawMode DrawMode, texture *Texture, shaderProgram *ShaderProgram, targetArea geometry.Rect) (*Sprite, error) {
	if texture == nil || texture.glHandler == 0 {
		return nil, fmt.Errorf("no texture")
	}

	if shaderProgram == nil || shaderProgram.glHandler == 0 {
		return nil, fmt.Errorf("no shader program")
	}

	vertices := geometry.ComputeSpriteVertices(width, height, targetArea)
	textureCoordinates := geometry.ComputeSpriteTextureCoordinates(
		texture.imageWidth, texture.imageHeight, targetArea)
	colorMask := RGBAFullOpaque()
	colorMaskData := colorMask.Data()

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
	gl.BufferData(gl.ARRAY_BUFFER, len(colorMaskData[:])*4, gl.Ptr(colorMaskData[:]), uint32(colorDrawMode))
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

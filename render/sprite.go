package render

import (
	_ "embed"
	"fmt"
	"sort"

	"github.com/alacrity-engine/core/geometry"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	spriteIndexBufferHandler uint32

	//go:embed std-sprite-vert.glsl
	standardSpriteVertexShaderSource string
	//go:embed std-sprite-frag.glsl
	standardSpriteFragmentShaderSource string
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
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 3*4, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glTextureCoordinatesBufferHandler)
	texCoordAttrib := uint32(gl.GetAttribLocation(sprite.shaderProgram.glHandler, gl.Str("aTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 2*4, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, sprite.glColorMaskBufferHandler)
	colorAttrib := uint32(gl.GetAttribLocation(sprite.shaderProgram.glHandler, gl.Str("aColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointerWithOffset(colorAttrib, 4, gl.FLOAT, false, 4*4, 0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (sprite *Sprite) SetZ(z float32) error {
	oldZ := sprite.drawZ
	sprite.drawZ = mgl32.Clamp(z, zMin, zMax)

	if sprite.batch != nil {
		// Remove sprite from the batch.
		ind := sprite.batchIndex
		batch := sprite.batch
		length := len(batch.sprites)

		// Reindex all the sprites
		// remaining on the batch.
		for i := ind + 1; i < length; i++ {
			batch.sprites[i].batchIndex--
		}

		if ind == 0 {
			batch.sprites = batch.sprites[1:]
		} else if ind < length-1 {
			batch.sprites = append(batch.sprites[:ind],
				batch.sprites[ind+1:]...)
		} else {
			batch.sprites = batch.sprites[:length-1]
		}

		err := batch.vertices.removeElements(ind, 18)

		if err != nil {
			return err
		}

		err = batch.texCoords.removeElements(ind, 12)

		if err != nil {
			return err
		}

		err = batch.colorMasks.removeElements(ind, 24)

		if err != nil {
			return err
		}

		err = batch.models.removeElements(ind, 16)

		if err != nil {
			return err
		}

		err = batch.shouldDraw.removeElement(ind)

		if err != nil {
			return err
		}

		// Attach sprite to the batch.
		ind = sort.Search(len(batch.sprites), func(i int) bool {
			return batch.sprites[i].drawZ >= sprite.drawZ
		})

		sprite.batchIndex = ind
		sprite.batch = batch

		if ind <= 0 {
			batch.sprites = append([]*Sprite{sprite}, batch.sprites...)
		} else if ind >= len(batch.sprites) {
			batch.sprites = append(batch.sprites, sprite)
		} else {
			batch.sprites = append(batch.sprites[:ind+1], batch.sprites[ind:]...)
			batch.sprites[ind] = sprite
		}

		// Reindex all the sprites
		// remaining on the batch.
		for i := ind + 1; i < len(batch.sprites); i++ {
			batch.sprites[i].batchIndex--
		}

		vertices := make([]float32, 18)
		geometry.ComputeSpriteVerticesNoElementsFill(
			vertices, width, height, sprite.targetArea)
		batch.vertices.insertElements(ind*len(vertices), len(vertices), vertices)

		texCoords := make([]float32, 12)
		geometry.ComputeSpriteTextureCoordinatesNoElementsFill(
			texCoords, sprite.texture.imageWidth,
			sprite.texture.imageHeight, sprite.targetArea)
		batch.texCoords.insertElements(ind*len(texCoords), len(texCoords), texCoords)

		colorMask := make([]float32, 24)
		geometry.ColorMaskDataNoElementsFill(
			colorMask, sprite.colorMask.Data())
		batch.colorMasks.insertElements(ind*len(colorMask), len(colorMask), colorMask)

		// Check the proportion.
		//a := batch.vertices.capacity / batch.vertices.stride
		//b := batch.texCoords.capacity / batch.texCoords.stride
		//c := batch.colorMasks.capacity / batch.colorMasks.stride
		//
		//if a != b && b != c && a != c {
		//	_ = a
		//}

		// Rebind all the sprite data to the VAO.
		batch.buildVAO()

		prevCapacity := batch.models.getCapacity()
		identMatrix := mgl32.Ident4()
		batch.models.insertElements(ind*len(identMatrix), len(identMatrix), identMatrix[:])

		if batch.models.getCapacity() > prevCapacity {
			batch.modelsTextureBuffer.
				rebind(batch.models.glHandler)
		}

		prevCapacity = batch.shouldDraw.getCapacity()
		batch.shouldDraw.insertElement(ind, 0)

		if batch.shouldDraw.getCapacity() > prevCapacity {
			batch.shouldDrawTextureBuffer.
				rebind(batch.shouldDraw.glHandler)
		}

		return nil
	}

	if sprite.canvas != nil {
		err := sprite.canvas.setSpriteZ(sprite, oldZ, sprite.drawZ)

		if err != nil {
			return err
		}
	}

	return nil
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
		sprite.batchIndex*len(colorMaskData), len(colorMaskData), colorMaskData)

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
	err := sprite.batch.texCoords.replaceElements(sprite.batchIndex*len(textureCoordinates),
		len(textureCoordinates), textureCoordinates)

	if err != nil {
		return err
	}

	vertices := make([]float32, 18)
	geometry.ComputeSpriteVerticesNoElementsFill(
		vertices, width, height, targetArea)
	err = sprite.batch.vertices.replaceElements(
		sprite.batchIndex*len(vertices), len(vertices), vertices)

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

	//gl.Disable(gl.DEPTH_TEST)
	//defer gl.Enable(gl.DEPTH_TEST)

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

	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)
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
		sprite.batchIndex*len(model), len(model), model[:])

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

	// Set the sprite to be drawn.
	sprite.canvas.sprites[sprite] = transform

	return nil
}

// TODO: encapsulate all the draw modes
// and shader program into sprite options.

func NewSpriteFromTextureAndProgram(vertexDrawMode, textureDrawMode, colorDrawMode DrawMode, texture *Texture, shaderProgram *ShaderProgram, targetArea geometry.Rect) (*Sprite, error) {
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
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), uint32(vertexDrawMode))
	vertAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 3*4, 0)

	var textureCoordinatesBufferHandler uint32
	gl.GenBuffers(1, &textureCoordinatesBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordinatesBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(textureCoordinates)*4, gl.Ptr(textureCoordinates), uint32(textureDrawMode))
	texCoordAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, 2*4, 0)

	var colorMaskBufferHandler uint32
	gl.GenBuffers(1, &colorMaskBufferHandler)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorMaskBufferHandler)
	gl.BufferData(gl.ARRAY_BUFFER, len(colorMaskData[:])*4, gl.Ptr(colorMaskData[:]), uint32(colorDrawMode))
	colorAttrib := uint32(gl.GetAttribLocation(shaderProgram.glHandler, gl.Str("aColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointerWithOffset(colorAttrib, 4, gl.FLOAT, false, 4*4, 0)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return &Sprite{
		glHandler:                         handler,
		glVertexBufferHandler:             vertexBufferHandler,
		glTextureCoordinatesBufferHandler: textureCoordinatesBufferHandler,
		glColorMaskBufferHandler:          colorMaskBufferHandler,
		texture:                           texture,
		targetArea:                        targetArea,
		colorMask:                         colorMask,
		shaderProgram:                     shaderProgram,
		drawMode:                          textureDrawMode,
		batchIndex:                        -1,
	}, nil
}

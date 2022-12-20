package render

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// TODO: collect views and projections
// assigned to cameras and canvases of
// the batched sprites. Pass a view
// and a projection index for each sprite.

// TODO: add a shader program field to batch.
// The projections and views should be
// uniform arrays of predefined sizes.
// The user can assign the initial size
// to the shader program, and if he needs
// to go beyond it, the shader program must
// be recompiled. The absolute max uniform
// array size cannot go past 256 because the
// projection index number for a vertice
// is only 1 byte long. Everytime the user
// adds a new canvas to the layout, the batch
// shader program uniforms are reassigned.

type Batch struct {
	glHandler uint32 // glHandler is an OpenGL name for the underlying batch VAO.
	sprites   []*Sprite
	layout    *Layout
	texture   *Texture

	// TODO: everytime we change a
	// parameter of the sprite we
	// should also change it in the
	// corresponding GPU list by the
	// batch index of the sprite.
	projectionsIdx *gpuList[byte]
	models         *gpuList[float32]
	viewsIdx       *gpuList[byte]
	vertices       *gpuList[float32]
	texCoords      *gpuList[float32]
	colorMasks     *gpuList[float32]
	shouldDraw     *gpuList[byte]
}

func (batch *Batch) Draw() {}

func (batch *Batch) AttachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	if sprite.texture != batch.texture {
		return fmt.Errorf(
			"the sprite should have the same texture as the batch")
	}

	ind := len(batch.sprites)
	sprite.batchIndex = ind
	sprite.batch = batch
	batch.sprites = append(batch.sprites, sprite)

	identMatrix := mgl32.Ident4()

	batch.vertices.addDataFromBuffer(
		sprite.glVertexBufferHandler, 8)
	batch.texCoords.addDataFromBuffer(
		sprite.glTextureCoordinatesBufferHandler, 8)
	batch.colorMasks.addDataFromBuffer(
		sprite.glColorMaskBufferHandler, 16)

	batch.projectionsIdx.addElement(byte(sprite.canvas.index))
	batch.viewsIdx.addElement(byte(sprite.canvas.index))
	batch.models.addElements(identMatrix[:])
	batch.shouldDraw.addElement(0)

	sprite.deleteVertexBuffer()
	sprite.deleteTextureCoordinatesBuffer()
	sprite.deleteColorMaskBuffer()
	sprite.deleteVertexArray()

	return nil
}

func (batch *Batch) DetachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	length := len(batch.sprites)
	ind := sprite.batchIndex

	if ind < length && batch.sprites[ind] != sprite ||
		ind >= length || ind < 0 {
		return fmt.Errorf(
			"the sprite doesn't exist on the batch")
	}

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

	sprite.batch = nil
	batchIndex := sprite.batchIndex
	sprite.batchIndex = -1

	sprite.createVertexBuffer()
	sprite.createTextureCoordinatesBuffer()
	sprite.createColorMaskBuffer()

	batch.vertices.copyDataToBuffer(sprite.
		glVertexBufferHandler, batchIndex, 8)
	batch.texCoords.copyDataToBuffer(sprite.
		glTextureCoordinatesBufferHandler, batchIndex, 8)
	batch.colorMasks.copyDataToBuffer(sprite.
		glColorMaskBufferHandler, batchIndex, 16)

	sprite.createVertexArray()
	sprite.assembleVertexArray()

	err := batch.vertices.removeElements(batchIndex, 8)

	if err != nil {
		return err
	}

	err = batch.texCoords.removeElements(batchIndex, 8)

	if err != nil {
		return err
	}

	err = batch.colorMasks.removeElements(batchIndex, 16)

	if err != nil {
		return err
	}

	err = batch.projectionsIdx.removeElement(batchIndex)

	if err != nil {
		return err
	}

	err = batch.viewsIdx.removeElement(batchIndex)

	if err != nil {
		return err
	}

	err = batch.models.removeElements(batchIndex, 16)

	if err != nil {
		return err
	}

	err = batch.shouldDraw.removeElement(batchIndex)

	if err != nil {
		return err
	}

	return nil
}

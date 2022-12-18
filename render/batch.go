package render

import (
	"fmt"
)

// TODO: collect views and projections
// assigned to cameras and canvases of
// the batched sprites. Pass a view
// and a projection index for each sprite.

// TODO: add a shader program field to batch.

// TODO: a vertex draw buffer for all the
// vertices of all the attached sprites.
// If a sprite shouldn't be drawn, set
// its color to 0 in all the batch shaders.

type Batch struct {
	glHandler uint32 // glHandler is an OpenGL name for the underlying batch VAO.
	sprites   []*Sprite
	layout    *Layout
	texture   *Texture

	// TODO: add buffers for the
	// data of all the sprites
	// on the batch (type: *gpuList).
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

	// TODO: copy all the sprite data
	// to the batch buffers.

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

	if ind == 0 {
		batch.sprites = batch.sprites[1:]
	} else if ind < length-1 {
		batch.sprites = append(batch.sprites[:ind],
			batch.sprites[ind+1:]...)
	} else {
		batch.sprites = batch.sprites[:length-1]
	}

	sprite.batch = nil
	sprite.batchIndex = -1

	// TODO: remove all the sprite
	// data from the batch buffers.

	return nil
}

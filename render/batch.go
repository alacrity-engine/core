package render

import (
	"fmt"
	"sort"
	"unsafe"

	"github.com/alacrity-engine/core/geometry"
)

// TODO: collect views and projections
// assigned to cameras and canvases of
// the batched sprites. Pass a view
// and a projection index for each sprite.

// TODO: add a shader program field to batch.

// TODO: a vertex draw buffer for all the
// vertices of all the attached sprites.
// If a sprite shouldn't be drawn, don't
// add its vertex data to the vertex
// buffer. When the vertex buffer is
// being filled with the vertex data of
// all the sprites that should be drawn,
// the quantity of the sprites that
// shouldn't be drawn must be counted in
// order to cut the length of the vertex
// buffer to the number of the sprites
// that should be drawn.

type Batch struct {
	// glHandler is an OpenGL name
	// for the underlying batch VAO.
	glHandler    uint32
	list         *glList[float32]
	sprites      []*Sprite
	transforms   []*geometry.Transform
	shouldDraw   []bool
	layout       *Layout
	texture      *Texture
	vertexBuffer []float32
}

func (batch *Batch) shouldCreateVertexBuffer() {
	if batch.vertexBuffer == nil {
		batch.vertexBuffer = make([]float32, len(batch.sprites)*20) // paste the actual size of the vertex here
	}
}

func (batch *Batch) Draw() {
	batch.shouldCreateVertexBuffer()
}

func (batch *Batch) findSpritePlaceIndex(sprite *Sprite) int {
	return sort.Search(len(batch.sprites), func(i int) bool {
		return uintptr(unsafe.Pointer(batch.sprites[i])) >=
			uintptr(unsafe.Pointer(sprite))
	})
}

func (batch *Batch) findSpriteIndex(sprite *Sprite) int {
	return sort.Search(len(batch.sprites), func(i int) bool {
		return uintptr(unsafe.Pointer(batch.sprites[i])) ==
			uintptr(unsafe.Pointer(sprite))
	})
}

func (batch *Batch) AttachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	if sprite.texture != batch.texture {
		return fmt.Errorf(
			"the sprite should have the same texture as the batch")
	}

	length := len(batch.sprites)
	ind := batch.findSpritePlaceIndex(sprite)

	if ind < length && batch.sprites[ind] == sprite {
		return fmt.Errorf("the sprite already exists on the batch")
	}

	if ind == 0 {
		batch.sprites = append(batch.sprites, nil)
		copy(batch.sprites[1:], batch.sprites)
		batch.sprites[0] = sprite
	} else if ind < length {
		batch.sprites = append(batch.sprites[:ind+1],
			batch.sprites[ind:]...)
		batch.sprites[ind] = sprite
	} else {
		batch.sprites = append(batch.sprites, sprite)
	}

	sprite.batch = batch
	batch.transforms = append(batch.transforms, nil)
	batch.shouldDraw = append(batch.shouldDraw, false)

	return nil
}

func (batch *Batch) DetachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	length := len(batch.sprites)
	ind := batch.findSpriteIndex(sprite)

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
	batch.transforms = batch.transforms[:length-1]
	batch.shouldDraw = batch.shouldDraw[:length-1]

	return nil
}

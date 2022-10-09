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

type Batch struct {
	glHandler  uint32
	length     int
	capacity   int
	sprites    []*Sprite
	transforms []*geometry.Transform
	texture    *Texture
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

	return nil
}

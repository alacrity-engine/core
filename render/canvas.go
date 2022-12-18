package render

import "fmt"

// TODO: add a projection field and
// use it in sprite.Draw() operations.

type Canvas struct {
	index   int
	sprites map[*Sprite]struct{}
	layout  *Layout
	camera  *Camera
}

func (canvas *Canvas) Camera() *Camera {
	return canvas.camera
}

func (canvas *Canvas) Index() int {
	return canvas.index
}

func (canvas *Canvas) Z() float32 {
	return 2.0 * float32(canvas.index)
}

func (canvas *Canvas) Range() (float32, float32) {
	return -1.0 + float32(canvas.index)*2.0,
		1.0 + float32(canvas.index)*2.0
}

func (canvas *Canvas) AddSprite(sprite *Sprite) error {
	if _, ok := canvas.sprites[sprite]; ok {
		return fmt.Errorf(
			"the sprite already exists on the canvas")
	}

	canvas.sprites[sprite] = struct{}{}
	sprite.canvas = canvas

	return nil
}

func (canvas *Canvas) RemoveSprite(sprite *Sprite) error {
	if _, ok := canvas.sprites[sprite]; !ok {
		return fmt.Errorf(
			"the sprite doesn't exist on the canvas")
	}

	delete(canvas.sprites, sprite)
	sprite.canvas = nil

	return nil
}

func NewCanvas(drawZ int) *Canvas {
	return &Canvas{
		sprites: map[*Sprite]struct{}{},
		index:   drawZ,
		camera:  NewCamera(),
	}
}

package render

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type Canvas struct {
	index      int
	pos        byte
	sprites    map[*Sprite]struct{}
	layout     *Layout
	camera     *Camera
	projection mgl32.Mat4
}

func (canvas *Canvas) Camera() *Camera {
	return canvas.camera
}

func (canvas *Canvas) Index() int {
	return canvas.index
}

func (canvas *Canvas) Z() float32 {
	zLength := zMax - zMin

	return zLength * float32(canvas.index)
}

func (canvas *Canvas) Range() (float32, float32) {
	zLength := zMax - zMin

	return zMin + float32(canvas.index)*zLength,
		zMax + float32(canvas.index)*zLength
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

func (canvas *Canvas) updateBatchViews() {
	for i := 0; i < len(canvas.layout.batches); i++ {
		batch := canvas.layout.batches[i]
		batch.setCanvasView(int(canvas.pos), canvas.camera.View())
	}
}

func NewCanvas(drawZ int, projection mgl32.Mat4) *Canvas {
	camera := NewCamera()
	canvas := &Canvas{
		sprites:    map[*Sprite]struct{}{},
		index:      drawZ,
		camera:     NewCamera(),
		projection: projection,
	}

	camera.canvas = canvas

	return canvas
}

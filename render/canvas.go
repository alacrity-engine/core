package render

import "github.com/go-gl/mathgl/mgl32"

type Canvas struct {
	sprites []*Sprite
	drawZ   int
}

func (canvas *Canvas) View() mgl32.Mat4 {
	return mgl32.Translate3D(0, 0, float32(canvas.drawZ))
}

func (canvas *Canvas) AddSprite(sprite *Sprite) {
	canvas.sprites = append(canvas.sprites, sprite)
	sprite.drawZ = canvas.drawZ
}

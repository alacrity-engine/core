package render

type Canvas struct {
	index   int
	sprites []*Sprite
	layout  *Layout
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

func (canvas *Canvas) AddSprite(sprite *Sprite) {
	canvas.sprites = append(canvas.sprites, sprite)
	sprite.canvas = canvas
}

func NewCanvas(drawZ int) *Canvas {
	return &Canvas{
		sprites: []*Sprite{},
		index:   drawZ,
	}
}

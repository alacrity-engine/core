package render

import "github.com/go-gl/gl/v4.6-core/gl"

type TextureBufferFormat uint32

const (
	FormatByte    TextureBufferFormat = gl.R8UI
	FormatFloat32 TextureBufferFormat = gl.R32F
)

type TextureBuffer struct {
	glHandler       uint32
	glBufferHandler uint32
	slot            TextureSlot
	format          TextureBufferFormat
}

func (tb *TextureBuffer) Bind() {
	gl.ActiveTexture(uint32(tb.slot))
	gl.BindTexture(gl.TEXTURE_BUFFER, tb.glHandler)
	gl.TexBuffer(gl.TEXTURE_BUFFER, uint32(tb.format), tb.glBufferHandler)

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)
}

func (tb *TextureBuffer) setBuffer(glBufferHandler uint32) {
	tb.glBufferHandler = glBufferHandler
}

func (tb *TextureBuffer) rebind(glBufferHandler uint32) {
	tb.setBuffer(glBufferHandler)
	tb.Bind()
}

func NewTextureBuffer(glBufferHandler uint32, slot TextureSlot, format TextureBufferFormat) *TextureBuffer {
	var glHandler uint32

	gl.GenTextures(1, &glHandler)
	gl.ActiveTexture(uint32(slot))
	gl.BindTexture(gl.TEXTURE_BUFFER, glHandler)
	gl.TexBuffer(gl.TEXTURE_BUFFER, uint32(format), glBufferHandler)

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)

	return &TextureBuffer{
		glHandler:       glHandler,
		glBufferHandler: glBufferHandler,
		slot:            slot,
		format:          format,
	}
}

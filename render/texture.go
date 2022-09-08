package render

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type TextureFiltering uint32

const (
	TextureFilteringNearest TextureFiltering = gl.NEAREST
	TextureFilteringLinear  TextureFiltering = gl.LINEAR
)

type Texture struct {
	glHandler   uint32
	imageWidth  int
	imageHeight int
	filter      TextureFiltering
}

func (texture *Texture) Use() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.glHandler)
}

func NewTextureFromImage(img *image.RGBA, filter TextureFiltering) *Texture {
	var handler uint32

	gl.GenTextures(1, &handler)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, handler)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(filter))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(filter))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Rect.Size().X),
		int32(img.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	gl.BindTexture(gl.TEXTURE_2D, 0)
	//gl.ActiveTexture(0)

	return &Texture{
		glHandler:   handler,
		imageWidth:  img.Rect.Max.X,
		imageHeight: img.Rect.Max.Y,
		filter:      filter,
	}
}

func NewEmptyTexture(width, height int, filter TextureFiltering) *Texture {
	var handler uint32

	gl.GenTextures(1, &handler)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, handler)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(filter))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(filter))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width),
		int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(nil))

	gl.BindTexture(gl.TEXTURE_2D, 0)
	//gl.ActiveTexture(0)

	return &Texture{
		glHandler:   handler,
		imageWidth:  width,
		imageHeight: height,
		filter:      filter,
	}
}

package render

import (
	"image/color"
	"reflect"
	"unsafe"
)

// RGBA represents an alpha-premultiplied RGBA color with components within range [0, 1].
//
// The difference between color.RGBA is that the value range is [0, 1] and the values are floats.
//
// Taken from github.com/faiface/pixel.
type RGBA struct {
	R, G, B, A float32
}

// RGB returns a fully opaque RGBA color with the given RGB values.
//
// A common way to construct a transparent color is to create one with RGB constructor, then
// multiply it by a color obtained from the Alpha constructor.
func RGB(r, g, b float32) RGBA {
	return RGBA{r, g, b, 1}
}

// Alpha returns a white RGBA color with the given alpha component.
func Alpha(a float32) RGBA {
	return RGBA{a, a, a, a}
}

// Add adds color d to color c component-wise and returns the result (the components are not
// clamped).
func (c RGBA) Add(d RGBA) RGBA {
	return RGBA{
		R: c.R + d.R,
		G: c.G + d.G,
		B: c.B + d.B,
		A: c.A + d.A,
	}
}

// Sub subtracts color d from color c component-wise and returns the result (the components
// are not clamped).
func (c RGBA) Sub(d RGBA) RGBA {
	return RGBA{
		R: c.R - d.R,
		G: c.G - d.G,
		B: c.B - d.B,
		A: c.A - d.A,
	}
}

// Mul multiplies color c by color d component-wise (the components are not clamped).
func (c RGBA) Mul(d RGBA) RGBA {
	return RGBA{
		R: c.R * d.R,
		G: c.G * d.G,
		B: c.B * d.B,
		A: c.A * d.A,
	}
}

// Scaled multiplies each component of color c by scale and returns the result (the components
// are not clamped).
func (c RGBA) Scaled(scale float32) RGBA {
	return RGBA{
		R: c.R * scale,
		G: c.G * scale,
		B: c.B * scale,
		A: c.A * scale,
	}
}

// RGBA returns alpha-premultiplied red, green, blue and alpha components of the RGBA color.
func (c RGBA) RGBA() (r, g, b, a uint32) {
	r = uint32(0xffff * c.R)
	g = uint32(0xffff * c.G)
	b = uint32(0xffff * c.B)
	a = uint32(0xffff * c.A)
	return
}

type ColorMask [4]RGBA

func (mask ColorMask) Data() [16]float32 {
	colorMaskSlice := mask[:]
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&colorMaskSlice))
	header.Len *= 4
	header.Cap *= 4
	colorMaskData := *(*[16]float32)(unsafe.Pointer(header.Data))

	return colorMaskData
}

func RGBARepeat4(_color RGBA) ColorMask {
	return ColorMask{
		_color,
		_color,
		_color,
		_color,
	}
}

func RGBAFullOpaque() ColorMask {
	opaque := RGBA{R: 1.0, G: 1.0, B: 1.0, A: 1.0}
	return RGBARepeat4(opaque)
}

// ToRGBA converts a color to RGBA format. Using this function is preferred to using RGBAModel, for
// performance (using RGBAModel introduces additional unnecessary allocations).
func ToRGBA(c color.Color) RGBA {
	if c, ok := c.(RGBA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA{
		float32(r) / 0xffff,
		float32(g) / 0xffff,
		float32(b) / 0xffff,
		float32(a) / 0xffff,
	}
}

// RGBAModel converts colors to RGBA format.
var RGBAModel = color.ModelFunc(rgbaModel)

func rgbaModel(c color.Color) color.Color {
	return ToRGBA(c)
}

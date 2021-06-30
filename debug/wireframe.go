package debug

import (
	"image/color"

	"github.com/alacrity-engine/core/system"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	colors "golang.org/x/image/colornames"
)

var (
	imd *imdraw.IMDraw
)

// WireframeInitialize initializes
// wireframe for drawing primitive
// forms (lines, circles, polygons).
func WireframeInitialize() {
	imd = imdraw.New(nil)
}

// WireframeClear clears
// the wireframe canvas.
func WireframeClear() {
	if imd == nil {
		return
	}

	imd.Clear()
}

// WireframeColor returns the
// color of the wireframe.
func WireframeColor() color.RGBA {
	if imd == nil {
		return colors.White
	}

	return imd.Color.(color.RGBA)
}

// WireframeSetColor assign a new color
// to the wireframe.
func WireframeSetColor(cl color.RGBA) {
	if imd == nil {
		return
	}

	imd.Color = cl
}

// WireframeDrawLine draws a straight line
// on the wireframe.
func WireframeDrawLine(a, b pixel.Vec, thickness float64) {
	if imd == nil {
		return
	}

	imd.Push(a)
	imd.Push(b)
	imd.Line(thickness)
}

// WireframeDrawCircle draws a new circle on the
// wireframe.
func WireframeDrawCircle(center pixel.Vec, radius, thickness float64) {
	if imd == nil {
		return
	}

	imd.Push(center)
	imd.Circle(radius, thickness)
}

// WireframeDrawPolygon draws a new polygon
// on the wireframe.
func WireframeDrawPolygon(thickness float64, points ...pixel.Vec) {
	if imd == nil {
		return
	}

	imd.Push(points...)
	imd.Polygon(thickness)
}

// WireframeDraw draws all
// the wireframe elements
// on the screen.
func WireframeDraw() {
	if imd == nil {
		return
	}

	imd.Draw(system.Window())
}

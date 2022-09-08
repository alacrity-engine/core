package convert

import (
	"github.com/faiface/pixel"
	"github.com/zergon321/cirno"
	"gonum.org/v1/plot/vg"
)

// GonumVGToPixelVector converts vg.Point (from
// gonum/plot) to pixel.Vec.
func GonumVGToPixelVector(point vg.Point) pixel.Vec {
	return pixel.V(float64(point.X), float64(point.Y))
}

// PixelVectorToGonumVG converts pixel.Vec to vg.Point
// (from gonum/plot).
func PixelVectorToGonumVG(vector pixel.Vec) vg.Point {
	return vg.Point{X: vg.Length(vector.X), Y: vg.Length(vector.Y)}
}

// VectorCirnoToPixel converts cirno.Vector to pixel.Vec.
func VectorCirnoToPixel(vector cirno.Vector) pixel.Vec {
	return pixel.V(vector.X, vector.Y)
}

// VectorPixelToCirno converts pixel.Vec to cirno.Vector.
func VectorPixelToCirno(vector pixel.Vec) cirno.Vector {
	return cirno.NewVector(vector.X, vector.Y)
}

// GonumVGToCirnoVector converts vg.Point (from gonum/plot)
// to cirno.Vector.
func GonumVGToCirnoVector(point vg.Point) cirno.Vector {
	return cirno.NewVector(float64(point.X), float64(point.Y))
}

// CirnoVectorToGonumVG converts cirno.Vector to vg.Point
// (from gonum/plot).
func CirnoVectorToGonumVG(vector cirno.Vector) vg.Point {
	return vg.Point{X: vg.Length(vector.X), Y: vg.Length(vector.Y)}
}

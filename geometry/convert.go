package geometry

import (
	"github.com/zergon321/cirno"
	"gonum.org/v1/plot/vg"
)

// GonumVGToAlacrityVector converts vg.Point (from
// gonum/plot) to pixel.Vec.
func GonumVGToAlacrityVector(point vg.Point) Vec {
	return V(float64(point.X), float64(point.Y))
}

// AlacrityVectorToGonumVG converts Vec to vg.Point
// (from gonum/plot).
func AlacrityVectorToGonumVG(vector Vec) vg.Point {
	return vg.Point{X: vg.Length(vector.X), Y: vg.Length(vector.Y)}
}

// VectorCirnoToAlacrity converts cirno.Vector to Vec.
func VectorCirnoToAlacrity(vector cirno.Vector) Vec {
	return V(vector.X, vector.Y)
}

// VectorAlacrityToCirno converts pixel.Vec to cirno.Vector.
func VectorAlacrityToCirno(vector Vec) cirno.Vector {
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

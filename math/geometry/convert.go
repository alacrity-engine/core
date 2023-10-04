package geometry

import (
	"math"

	"github.com/zergon321/cirno"
	"gonum.org/v1/plot/vg"
)

const (
	// RadToDeg is a factor to transfrom radians to degrees.
	RadToDeg float64 = 180.0 / math.Pi
	// DegToRad is a factor to transform degrees to radians.
	DegToRad float64 = math.Pi / 180.0
)

// AdjustAngle adjusts the value of the angle so it
// is bettween 0 and 360.
func AdjustAngle(angle float64) float64 {
	// Adjust the angle so its value is between 0 and 360.
	if angle >= 360 {
		angle = angle - float64(int64(angle/360))*360
	} else if angle < 0 {
		if angle <= -360 {
			angle = angle - float64(int64(angle/360))*360
		}

		angle += 360

		if angle >= 360 {
			angle = angle - float64(int64(angle/360))*360
		}
	}

	return angle
}

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

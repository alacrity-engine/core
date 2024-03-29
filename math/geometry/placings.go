package geometry

import (
	"math"

	"gonum.org/v1/plot/tools/bezier"
	"gonum.org/v1/plot/vg"
)

// generatePoints generate points of the trajectory
// using parametrized functions for X and Y coordinates.
func generatePoints(x, y func(float64) float64, numberOfPoints int) []Vec {
	points := []Vec{}

	for i := 0; i < numberOfPoints; i++ {
		t := float64(i) / float64(numberOfPoints)
		point := V(x(t), y(t))

		points = append(points, point)
	}

	return points
}

// BezierCurve returns a new Bezier curve consisting of
// equidistant points.
func PlacingBezierCurve(controlPoints []Vec, numberOfPoints int, dt float64) []Vec {
	vgControlPoints := []vg.Point{}

	for _, controlPoint := range controlPoints {
		vgControlPoint := AlacrityVectorToGonumVG(controlPoint)
		vgControlPoints = append(vgControlPoints, vgControlPoint)
	}

	curve := bezier.New(vgControlPoints...)
	vgPoints := []vg.Point{}

	for t := 0.0; t < 100; t += dt {
		point := curve.Point(t / 100.0)
		vgPoints = append(vgPoints, point)
	}

	points := []Vec{}

	for _, vgPoint := range vgPoints {
		point := GonumVGToAlacrityVector(vgPoint)
		points = append(points, point)
	}

	curvePoints := GetSegmentPoints(points,
		numberOfPoints-1)

	return curvePoints
}

// TODO: add more Jordan curves;
// rewrite Line() with generatePoints().

// Line returns a new line consisting of equidistant points.
func PlacingLine(start, end Vec, numberOfPoints int) []Vec {
	linePoints := []Vec{}
	step := end.Sub(start).Scaled(1.0 / (float64(numberOfPoints) - 1))

	linePoints = append(linePoints, start)

	for i := 0; i < numberOfPoints-1; i++ {
		point := start.Add(step.Scaled(float64(i + 1)))
		linePoints = append(linePoints, point)
	}

	return linePoints
}

// Ellipse creates a new ellipse consisting of points.
func PlacingEllipse(center Vec, a, b float64, numberOfPoints int) []Vec {
	x := func(t float64) float64 {
		return a*math.Cos(2*math.Pi*t) + center.X
	}
	y := func(t float64) float64 {
		return b*math.Sin(2*math.Pi*t) + center.Y
	}

	rawPoints := generatePoints(x, y, numberOfPoints)
	equidistantPoints := GetSegmentPoints(rawPoints, numberOfPoints-1)

	return equidistantPoints
}

// Circle creates a new circle consisting of points.
func PlacingCircle(center Vec, radius float64, numberOfPoints int) []Vec {
	return PlacingEllipse(center, radius, radius, numberOfPoints)
}

// Astroid creates a new astroid consisting of points.
func PlacingAstroid(center Vec, a, b float64, numberOfPoints int) []Vec {
	x := func(t float64) float64 {
		return a*math.Pow(math.Cos(2*math.Pi*t), 3) + center.X
	}
	y := func(t float64) float64 {
		return b*math.Pow(math.Sin(2*math.Pi*t), 3) + center.Y
	}

	return generatePoints(x, y, numberOfPoints)
}

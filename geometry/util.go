package geometry

import (
	"math"
)

const (
	// Epsilon is a number for
	// float comparison.
	Epsilon = 0.001
)

// GetSegmentPoints returns the points of
// the equally segmented curve.
func GetSegmentPoints(points []Vec, numberOfSegments int) []Vec {
	// Create lines out of bezier
	// curve points.
	lines := []Line{}

	for i := 0; i < len(points)-1; i++ {
		line := L(points[i],
			points[i+1])

		lines = append(lines, line)
	}

	// Compute the length
	// of the bezier curve
	// interpolated with lines.
	length := 0.0

	for _, line := range lines {
		length += line.Len()
	}

	// Divide the bezier curve into
	// equal segments.
	step := length / float64(numberOfSegments)
	segmentPoints := []Vec{}
	lastLine := 0
	lastPoint := lines[0].A
	segmentPoints = append(segmentPoints, lastPoint)

	for i := 0; i < numberOfSegments; i++ {
		startLine := L(lastPoint, lines[lastLine].B)
		localLength := startLine.Len()

		for step-localLength > Epsilon {
			line := lines[lastLine+1]

			localLength += line.Len()
			lastLine++
		}

		line := lines[lastLine]

		if localLength-step > Epsilon {
			difference := localLength - step
			t := difference / line.Len()

			lastPoint = V(t*line.A.X+(1-t)*line.B.X,
				t*line.A.Y+(1-t)*line.B.Y)
		} else {
			lastPoint = line.B
			lastLine++
		}

		segmentPoints = append(segmentPoints, lastPoint)
	}

	return segmentPoints
}

// ClampMagnitude clamps the vector to the specified magnitude
// if its length exceeds it.
func ClampMagnitude(vector Vec, magnitude float64) Vec {
	currentMagnitude := vector.Len()
	factor := math.Min(currentMagnitude, magnitude) / currentMagnitude

	return vector.Scaled(factor)
}

// Clamp returns x clamped to the interval [min, max].
//
// If x is less than min, min is returned. If x is more than max, max is returned. Otherwise, x is
// returned.
//
// Taken from https://github.com/faiface/pixel
func Clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}

	if x > max {
		return max
	}

	return x
}

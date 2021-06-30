package geometry

import (
	"math"

	"github.com/faiface/pixel"
)

const (
	// Epsilon is a number for
	// float comparison.
	Epsilon = 0.001
)

// GetSegmentPoints returns the points of
// the equally segmented curve.
func GetSegmentPoints(points []pixel.Vec, numberOfSegments int) []pixel.Vec {
	// Create lines out of bezier
	// curve points.
	lines := []pixel.Line{}

	for i := 0; i < len(points)-1; i++ {
		line := pixel.L(points[i],
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
	segmentPoints := []pixel.Vec{}
	lastLine := 0
	lastPoint := lines[0].A
	segmentPoints = append(segmentPoints, lastPoint)

	for i := 0; i < numberOfSegments; i++ {
		subsegments := []pixel.Line{}
		startLine := pixel.L(lastPoint, lines[lastLine].B)

		subsegments = append(subsegments, startLine)
		localLength := startLine.Len()

		for step-localLength > Epsilon {
			line := lines[lastLine+1]
			subsegments = append(subsegments, line)

			localLength += line.Len()
			lastLine++
		}

		line := lines[lastLine]

		if localLength-step > Epsilon {
			difference := localLength - step
			t := difference / line.Len()

			lastPoint = pixel.V(t*line.A.X+(1-t)*line.B.X,
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
func ClampMagnitude(vector pixel.Vec, magnitude float64) pixel.Vec {
	currentMagnitude := vector.Len()
	factor := math.Min(currentMagnitude, magnitude) / currentMagnitude

	return vector.Scaled(factor)
}

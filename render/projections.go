package render

import "github.com/go-gl/mathgl/mgl32"

var (
	ortho2D mgl32.Mat4
)

func ortho2DCompute() mgl32.Mat4 {
	aspect := float32(height) / float32(width)
	return mgl32.Ortho(-1, 1, -1*aspect, 1*aspect, -1, 1)
}

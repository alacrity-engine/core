package render

import "github.com/go-gl/mathgl/mgl32"

func Ortho2DStandard() mgl32.Mat4 {
	aspect := float32(height) / float32(width)
	return mgl32.Ortho(-1, 1, -1*aspect, 1*aspect, zMin, zMax)
}

func Ortho2D(zNear, zFar float32) mgl32.Mat4 {
	aspect := float32(height) / float32(width)
	return mgl32.Ortho(-1, 1, -1*aspect, 1*aspect, zNear, zFar)
}

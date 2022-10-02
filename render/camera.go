package render

import (
	"github.com/alacrity-engine/core/geometry"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	transform *geometry.Transform
}

func (camera *Camera) View() mgl32.Mat4 {
	return camera.transform.Data()
}

func (camera *Camera) Move(direction geometry.Vec) *Camera {
	camera.transform.Move(direction)
	return camera
}

func (camera *Camera) MoveTo(destination geometry.Vec) *Camera {
	camera.transform.MoveTo(destination)
	return camera
}

func NewCamera() *Camera {
	return &Camera{
		transform: geometry.NewTransform(nil),
	}
}

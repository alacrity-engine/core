package geometry

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang-collections/collections/queue"
)

// Transform stores position, angle
// and scale of the game object.
type Transform struct {
	parent   *Transform
	model    mgl32.Mat4
	position Vec
	angle    float64
	scale    Vec
	children []*Transform
	z        float32
}

func (t *Transform) AdjustZ(z float32) *Transform {
	t.z += z
	t.model = mgl32.Translate3D(0,
		0, z).Mul4(t.model)

	return t
}

// findChild returns the index of the given transform
// in the array of children.
func (t *Transform) findChild(tr *Transform) int {
	ind := -1

	for i, child := range t.children {
		if child == tr {
			ind = i
			break
		}
	}

	return ind
}

// Parent returns the parent of the transform.
func (t *Transform) Parent() *Transform {
	return t.parent
}

// SetParent sets the parent for the transform.
func (t *Transform) SetParent(parent *Transform) error {
	if parent == nil {
		if t.parent != nil {
			t.parent.RemoveChild(t)
		}

		return nil
	}

	return parent.AddChild(t)
}

// HasChild returns true if the transform has 'tr'
// as a direct child.
func (t *Transform) HasChild(tr *Transform) bool {
	i := t.findChild(tr)

	return i >= 0
}

// AddChild adds a new child to the transform.
func (t *Transform) AddChild(child *Transform) error {
	if t.HasChild(child) {
		return fmt.Errorf("the transform already has child '%v'",
			child)
	}

	t.children = append(t.children, child)
	child.parent = t

	return nil
}

// RemoveChild removes the child from the transform.
func (t *Transform) RemoveChild(child *Transform) error {
	i := t.findChild(child)

	if i < 0 {
		return fmt.Errorf("the transform has no child '%v'",
			child)
	}

	t.children = append(t.children[:i],
		t.children[i+1:]...)
	child.parent = nil

	return nil
}

// Axes returns the local axes (X and Y) of the transform.
func (t *Transform) Axes() (Vec, Vec) {
	sinA := math.Sin(t.angle)
	cosA := 1 - sinA*sinA

	return V(cosA, sinA), V(-sinA, cosA)
}

// Position returns the current position
// stored in the transform.
func (t *Transform) Position() Vec {
	return t.position
}

// Angle returns the transform angle in radians.
func (t *Transform) Angle() float64 {
	return t.angle
}

// Scale returns the scale of the transform.
func (t *Transform) Scale() Vec {
	return t.scale
}

// Data returns the matrix held by the
// given transform.
func (t *Transform) Data() mgl32.Mat4 {
	return t.model
}

// Move moves the transform and its children
// in the specified direction.
func (t *Transform) Move(direction Vec) *Transform {
	t.position = t.position.Add(direction)
	t.model = mgl32.Translate3D(float32(direction.X),
		float32(direction.Y), 0).Mul4(t.model)

	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)

		transform.position = transform.position.Add(direction)
		transform.model = mgl32.Translate3D(float32(direction.X),
			float32(direction.Y), 0).Mul4(transform.model)

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// MoveTo computes the translation by the current
// position and the destination and then performs
// movement of the transform and all its children.
func (t *Transform) MoveTo(position Vec) *Transform {
	offset := position.Sub(t.Position())
	return t.Move(offset)
}

// Rotate rotates the transform and all its children
// at the specified angle in degrees.
func (t *Transform) Rotate(angle float64) *Transform {
	t.angle += angle
	t.angle = AdjustAngle(t.angle)

	t.model = t.model.Mul4(mgl32.HomogRotate3DZ(
		float32(angle * DegToRad)))

	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)

		// Rotate the child transform itself.
		transform.angle += angle
		transform.angle = AdjustAngle(transform.angle)

		transform.model = transform.model.Mul4(mgl32.HomogRotate3DZ(
			float32(angle * DegToRad)))

		// Rotate the child transform
		// around the parent transform.
		base := transform.parent.Position()
		destination := transform.position.RotatedAround(angle*DegToRad, base)
		offset := destination.Sub(transform.Position())

		transform.position = transform.position.Add(offset)
		transform.model = mgl32.Translate3D(float32(offset.X),
			float32(offset.Y), 0).Mul4(transform.model)

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// RotateAround rotates the transform and all its
// children at the specified angle in degrees around
// the given point.
func (t *Transform) RotateAround(angle float64, base Vec) *Transform {
	destination := t.position.RotatedAround(angle*DegToRad, base)
	return t.MoveTo(destination)
}

// ApplyScale applies the specified scale factor
// to the transform and all its children.
func (t *Transform) ApplyScale(factor Vec) *Transform {
	t.scale = t.scale.Add(factor)
	t.model = t.model.Mul4(mgl32.Scale3D(
		float32(factor.X), float32(factor.Y), 0))

	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)

		transform.scale = transform.scale.Add(factor)
		transform.model = transform.model.Mul4(mgl32.Scale3D(
			float32(factor.X), float32(factor.Y), 0))

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// NewTransform creates a new empty transform out of given data.
func NewTransform(parent *Transform) *Transform {
	return &Transform{
		parent:   parent,
		model:    mgl32.Ident4(),
		children: []*Transform{},
	}
}

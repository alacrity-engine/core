package geometry

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/golang-collections/collections/queue"
)

// Transform stores position, angle
// and scale of the game object.
type Transform struct {
	parent   *Transform
	data     pixel.Matrix
	children []*Transform
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
func (t *Transform) Axes() (pixel.Vec, pixel.Vec) {
	tX := pixel.V(t.data[0], t.data[1]).Unit()
	tY := pixel.V(t.data[2], t.data[3]).Unit()

	return tX, tY
}

// Position returns the current position
// stored in the transform.
func (t *Transform) Position() pixel.Vec {
	return pixel.V(t.data[4], t.data[5])
}

// Angle returns the transform angle in radians.
func (t *Transform) Angle() float64 {
	localXaxis := pixel.V(t.data[0], t.data[1])

	return localXaxis.Angle()
}

// Scale returns the scale of the transform.
func (t *Transform) Scale() pixel.Vec {
	xAxis := pixel.V(t.data[0], t.data[1])
	yAxis := pixel.V(t.data[2], t.data[3])
	scale := pixel.V(xAxis.Len(), yAxis.Len())

	return scale
}

// Data returns the matrix held by the
// given transform.
func (t *Transform) Data() pixel.Matrix {
	return t.data
}

// Move moves the transform and its children
// in the specified direction.
func (t *Transform) Move(direction pixel.Vec) *Transform {
	t.data = t.data.Moved(direction)
	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)
		transform.data = transform.data.Moved(direction)

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// MoveTo computes the translation by the current
// position and the destination and then performs
// movement of the transform and all its children.
func (t *Transform) MoveTo(position pixel.Vec) *Transform {
	offset := position.Sub(t.Position())

	return t.Move(offset)
}

// Rotate rotates the transform and all its children
// at the specified angle.
func (t *Transform) Rotate(angle float64) *Transform {
	t.data = t.data.Rotated(t.Position(), angle)
	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)

		transform.data = transform.data.Rotated(
			transform.Position(), angle)
		transform.data = transform.data.Rotated(
			transform.parent.Position(), angle)

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// RotateAround rotates the transform and all its
// children at the specified angle in radians around
// the given point.
func (t *Transform) RotateAround(angle float64, base pixel.Vec) *Transform {
	t.data = t.data.Rotated(base, angle)
	//t.data = t.data.Rotated(t.Position(), angle)
	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)

		transform.data = transform.data.Rotated(
			base, angle)
		transform.data = transform.data.Rotated(
			transform.Position(), angle)

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// ApplyScale applies the specified scale factor
// to the transform and all its children.
func (t *Transform) ApplyScale(factor pixel.Vec) *Transform {
	t.data = t.data.ScaledXY(t.Position(), factor)
	tQueue := queue.New()

	for _, transform := range t.children {
		tQueue.Enqueue(transform)
	}

	for tQueue.Len() > 0 {
		transform := tQueue.Dequeue().(*Transform)

		transform.data = transform.data.ScaledXY(
			transform.Position(), factor)

		for _, child := range transform.children {
			tQueue.Enqueue(child)
		}
	}

	return t
}

// NewTransform creates a new empty transform out of given data.
func NewTransform(parent *Transform, data pixel.Matrix) *Transform {
	return &Transform{
		parent:   parent,
		data:     data,
		children: []*Transform{},
	}
}

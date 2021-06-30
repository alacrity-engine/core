package draw

import (
	"fmt"
	"image/color"

	"github.com/alacrity-engine/core/ecs"
	"github.com/alacrity-engine/core/system"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// Layout is a set of targets
// to draw game objects onto.
type Layout struct {
	groups []*drawGroup
}

// findGroup finds and returns the group with the specified
// name and its index, or nil if it cannot find the group.
func (layout *Layout) findGroup(name string) (int, *drawGroup) {
	ind := -1
	var group *drawGroup

	for i, drawGroup := range layout.groups {
		if name == drawGroup.name {
			ind = i
			group = drawGroup

			break
		}
	}

	return ind, group
}

// HasTarget indicates if the target with the specified
// name is on the layout.
func (layout *Layout) HasTarget(name string) bool {
	_, group := layout.findGroup(name)

	return group != nil
}

// AddTarget adds a new target onto the layout.
func (layout *Layout) AddTarget(name string, priority int, target pixel.Target) error {
	if layout.HasTarget(name) {
		return fmt.Errorf("draw froup with name '%s' already exists",
			name)
	}

	switch target.(type) {
	case *pixel.Batch, *pixelgl.Canvas:

	default:
		return fmt.Errorf("target must be either canvas or batch")
	}

	length := len(layout.groups)
	group := &drawGroup{
		name:      name,
		target:    target,
		gmobs:     []*ecs.GameObject{},
		transform: pixel.IM,
	}

	if length <= 0 || priority >= length {
		layout.groups = append(layout.groups, group)
	} else if priority < 0 {
		layout.groups = append([]*drawGroup{group},
			layout.groups...)
	} else {
		layout.groups = append(layout.groups[:priority+1],
			layout.groups[priority:]...)
		layout.groups[priority] = group
	}

	return nil
}

// PasteTargetBefore inserts the given target
// before the target that already exists.
func (layout *Layout) PasteTargetBefore(name, beforeName string, target pixel.Target) error {
	if layout.HasTarget(name) {
		return fmt.Errorf("draw froup with name '%s' already exists",
			name)
	}

	switch target.(type) {
	case *pixel.Batch, *pixelgl.Canvas:

	default:
		return fmt.Errorf("target must be either canvas or batch")
	}

	i, beforeGroup := layout.findGroup(beforeName)

	if beforeGroup == nil {
		return fmt.Errorf("there is no target with name '%s'",
			beforeName)
	}

	i--

	return layout.AddTarget(name, i, target)
}

// PasteTargetAfter pastes the the given target
// after the target that alewady exists.
func (layout *Layout) PasteTargetAfter(name, afterName string, target pixel.Target) error {
	if layout.HasTarget(name) {
		return fmt.Errorf("draw froup with name '%s' already exists",
			name)
	}

	switch target.(type) {
	case *pixel.Batch, *pixelgl.Canvas:

	default:
		return fmt.Errorf("target must be either canvas or batch")
	}

	i, afterGroup := layout.findGroup(afterName)

	if afterGroup == nil {
		return fmt.Errorf("there is no target with name '%s'",
			afterName)
	}

	i++

	return layout.AddTarget(name, i, target)
}

// TargetTransform returns the transform of the target.
func (layout *Layout) TargetTransform(name string) (pixel.Matrix, error) {
	_, group := layout.findGroup(name)

	if group == nil {
		return pixel.IM, fmt.Errorf("there is no target with name '%s'",
			name)
	}

	return group.transform, nil
}

// SetTargetTransform sets the transform to the target
// to affect how it's drawn on the window (used only
// with pixelgl.Canvas).
func (layout *Layout) SetTargetTransform(name string, transform pixel.Matrix) error {
	_, group := layout.findGroup(name)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			name)
	}

	group.transform = transform

	return nil
}

// AddGameObjectToTarget adds a new game object onto the target
// with the specified name.
func (layout *Layout) AddGameObjectToTarget(targetName string, gmob *ecs.GameObject, priority int) error {
	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	err := group.addGameObject(gmob, priority)

	return err
}

// RemoveTarget removes the target with the specified name
// from the layout.
func (layout *Layout) RemoveTarget(name string) error {
	i, group := layout.findGroup(name)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			name)
	}

	layout.groups = append(layout.groups[:i], layout.groups[i+1:]...)

	return nil
}

// RemoveGameObjectFromTarget removes the game object
// from the specified layout.
func (layout *Layout) RemoveGameObjectFromTarget(targetName, gmobName string) error {
	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	err := group.removeGameObject(gmobName)

	return err
}

// TargetCount returns the number of targets
// in the layout.
func (layout *Layout) TargetCount() int {
	return len(layout.groups)
}

// RenderOnTarget sends the sprite to be once
// rendered on the specified target and then removed.
func (layout *Layout) RenderOnTarget(targetName string, sprite *pixel.Sprite, transform pixel.Matrix, mask color.Color) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	group.addRenderObject(sprite, transform, mask)

	return nil
}

// RenderTextOnTarget sends the text to be rendered
// once on the specified target using the transform
// and the color mask.
func (layout *Layout) RenderTextOnTarget(targetName string, txt *text.Text, transform pixel.Matrix, mask color.Color) error {
	if txt == nil {
		return fmt.Errorf("the text is nil")
	}

	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	group.addTextObject(txt, transform, mask)

	return nil
}

// Draw all the game objects onto
// their targets and then deaws all
// the targets onto the game window.
func (layout *Layout) Draw() {
	for _, group := range layout.groups {
		// Remove all the destroyed game objects
		// before drawing.
		group.removeDestroyedGameObjects()

		// Clear the target.
		switch t := group.target.(type) {
		case *pixel.Batch:
			t.Clear()

		case *pixelgl.Canvas:
			t.Clear(color.RGBA{R: 0, G: 0, B: 0, A: 0})
		}

		// Draw game objects onto the target.
		for _, gmob := range group.gmobs {
			gmob.Draw(group.target)
		}

		// Draw all the render objects
		// and remove them.
		for _, rnob := range group.rnobs {
			rnob.sprite.DrawColorMask(group.target,
				rnob.transform, rnob.mask)
		}

		group.removeRenderObjects()

		// Draw all the text objects
		// and remove them.
		for _, txob := range group.txobs {
			txob.txt.DrawColorMask(group.target,
				txob.transform, txob.mask)
			txob.txt.Clear()
		}

		group.removeTextObjects()

		// Draw the target onto the window.
		switch t := group.target.(type) {
		case *pixel.Batch:
			t.Draw(system.Window())

		case *pixelgl.Canvas:
			t.Draw(system.Window(), group.transform)
		}
	}
}

// NewLayout returns a new
// layout to draw game objects
// onto targets.
func NewLayout() *Layout {
	return &Layout{
		groups: []*drawGroup{},
	}
}

package engine

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/alacrity-engine/core/system"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// DrawLayout is a set of targets
// to draw game objects onto.
type DrawLayout struct {
	groups []*drawGroup
}

// findGroup finds and returns the group with the specified
// name and its index, or nil if it cannot find the group.
func (layout *DrawLayout) findGroup(name string) (int, *drawGroup) {
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
func (layout *DrawLayout) HasTarget(name string) bool {
	_, group := layout.findGroup(name)

	return group != nil
}

// AddTarget adds a new target onto the layout.
func (layout *DrawLayout) AddTarget(name string, zDraw float64, target pixel.Target) error {
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
		gmobs:     []*GameObject{},
		transform: pixel.IM,
	}
	ind := sort.Search(length, func(i int) bool {
		return layout.groups[i].zDraw >= zDraw
	})

	if ind == 0 {
		layout.groups = append([]*drawGroup{group},
			layout.groups...)
	} else if ind < length {
		layout.groups = append(layout.groups[:ind+1],
			layout.groups[ind:]...)
		layout.groups[ind] = group
	} else {
		layout.groups = append(layout.groups, group)
	}

	group.zDraw = zDraw

	return nil
}

// TargetTransform returns the transform of the target.
func (layout *DrawLayout) TargetTransform(name string) (pixel.Matrix, error) {
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
func (layout *DrawLayout) SetTargetTransform(name string, transform pixel.Matrix) error {
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
func (layout *DrawLayout) AddGameObjectToTarget(targetName string, gmob *GameObject, zDraw float64) error {
	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	err := group.addGameObject(gmob, zDraw)

	return err
}

// RemoveTarget removes the target with the specified name
// from the layout.
func (layout *DrawLayout) RemoveTarget(name string) error {
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
func (layout *DrawLayout) RemoveGameObjectFromTarget(targetName, gmobName string) error {
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
func (layout *DrawLayout) TargetCount() int {
	return len(layout.groups)
}

// RenderOnTarget sends the sprite to be once
// rendered on the specified target and then removed.
func (layout *DrawLayout) RenderOnTarget(targetName string, sprite *pixel.Sprite, transform pixel.Matrix, mask color.Color, zDraw float64) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	group.addRenderObject(sprite, transform, mask, zDraw)

	return nil
}

// RenderTextOnTarget sends the text to be rendered
// once on the specified target using the transform
// and the color mask.
func (layout *DrawLayout) RenderTextOnTarget(targetName string, txt *text.Text, transform pixel.Matrix, mask color.Color, zDraw float64) error {
	if txt == nil {
		return fmt.Errorf("the text is nil")
	}

	_, group := layout.findGroup(targetName)

	if group == nil {
		return fmt.Errorf("there is no target with name '%s'",
			targetName)
	}

	group.addTextObject(txt, transform, mask, zDraw)

	return nil
}

// Draw all the game objects onto
// their targets and then deaws all
// the targets onto the game window.
func (layout *DrawLayout) Draw() {
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
func NewLayout() *DrawLayout {
	return &DrawLayout{
		groups: []*drawGroup{},
	}
}

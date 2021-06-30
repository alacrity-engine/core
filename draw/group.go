package draw

import (
	"fmt"
	"image/color"

	"github.com/alacrity-engine/core/ecs"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

// drawGroup is a group of
// game objects to be drawn
// on the same target.
type drawGroup struct {
	name      string
	target    pixel.Target
	gmobs     []*ecs.GameObject
	rnobs     []renderObject
	txobs     []textObject
	transform pixel.Matrix
}

// renderObject is a sprite with
// transform to be once rendered
// and then removed.
type renderObject struct {
	transform pixel.Matrix
	sprite    *pixel.Sprite
	mask      color.Color
}

// textObject is a text with
// transform to be once rendered
// and then removed.
type textObject struct {
	transform pixel.Matrix
	txt       *text.Text
	mask      color.Color
}

// hasGameObject returns true if the draw group has
// a game object with the specified name, and false
// otherwise.
func (group *drawGroup) hasGameObject(name string) bool {
	_, gmob := group.findGameObject(name)

	return gmob != nil
}

// addGameObject adds a new game object in the
// draw group computing its position by the priority.
func (group *drawGroup) addGameObject(gmob *ecs.GameObject, priority int) error {
	if group.hasGameObject(gmob.Name()) {
		return fmt.Errorf("the game object with name'%s'"+
			" already exists in the draw group", gmob.Name())
	}

	length := len(group.gmobs)

	if length <= 0 || priority >= length {
		group.gmobs = append(group.gmobs, gmob)
	} else if priority < 0 {
		temp := group.gmobs[1:]

		group.gmobs = []*ecs.GameObject{gmob}
		group.gmobs = append(group.gmobs, temp...)
	} else {
		group.gmobs = append(group.gmobs[:priority+1],
			group.gmobs[priority:]...)
		group.gmobs[priority] = gmob
	}

	return nil
}

// findGameObject returns the game object with the specified name
// or nil if t's not in the group.
func (group *drawGroup) findGameObject(name string) (int, *ecs.GameObject) {
	ind := -1
	var gameObject *ecs.GameObject

	for i, gmob := range group.gmobs {
		if name == gmob.Name() {
			ind = i
			gameObject = gmob

			break
		}
	}

	return ind, gameObject
}

// removeGameObject removes the game object with
// the specified name from the draw group.
func (group *drawGroup) removeGameObject(name string) error {
	i, gmob := group.findGameObject(name)

	if gmob == nil {
		return fmt.Errorf("there is no game object with name'%s'",
			name)
	}

	group.gmobs = append(group.gmobs[:i], group.gmobs[i+1:]...)

	return nil
}

// gameObjectCount returns the number of game
// objects placed in the draw group.
func (group *drawGroup) gameObjectCount() int {
	return len(group.gmobs)
}

// removeDestroyedGameObjects removes all the
// destroyed game objects from the draw group.
func (group *drawGroup) removeDestroyedGameObjects() {
	gmobs := []*ecs.GameObject{}

	// Find destroyed game objects.
	for _, gmob := range group.gmobs {
		if !gmob.Destroyed() {
			gmobs = append(gmobs, gmob)
		}
	}

	group.gmobs = gmobs
}

// addRenderObject adds a render object to the draw group
// so it will be drawn in the current frame and removed right after it.
func (group *drawGroup) addRenderObject(sprite *pixel.Sprite, transform pixel.Matrix, mask color.Color) {
	group.rnobs = append(group.rnobs,
		renderObject{sprite: sprite, transform: transform, mask: mask})
}

// addTextObject adds the text object to the draw group
// so it will be drawn in the current frame and removed right after it.
func (group *drawGroup) addTextObject(txt *text.Text, transform pixel.Matrix, mask color.Color) {
	group.txobs = append(group.txobs,
		textObject{txt: txt, transform: transform, mask: mask})
}

// removeRenderObjects removes all the render objects
// from the buffer.
func (group *drawGroup) removeRenderObjects() {
	group.rnobs = []renderObject{}
}

// removeTextObjects removes all the text objects
// from the buffer.
func (group *drawGroup) removeTextObjects() {
	group.txobs = []textObject{}
}

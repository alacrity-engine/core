package ecs

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

// drawGroup is a group of
// game objects to be drawn
// on the same target.
type drawGroup struct {
	name      string
	target    pixel.Target
	gmobs     []*GameObject
	rnobs     []renderObject
	txobs     []textObject
	transform pixel.Matrix
	zDraw     float64
}

// renderObject is a sprite with
// transform to be once rendered
// and then removed.
type renderObject struct {
	transform pixel.Matrix
	sprite    *pixel.Sprite
	mask      color.Color
	zDraw     float64
}

// textObject is a text with
// transform to be once rendered
// and then removed.
type textObject struct {
	transform pixel.Matrix
	txt       *text.Text
	mask      color.Color
	zDraw     float64
}

// hasGameObject returns true if the draw group has
// a game object with the specified name, and false
// otherwise.
func (group *drawGroup) hasGameObject(name string) bool {
	_, gmob := group.findGameObject(name)

	return gmob != nil
}

// insertGameObject inserts the game object
// into the sorted Z-buffer using binary search.
func (group *drawGroup) insertGameObject(gmob *GameObject, zDraw float64) {
	length := len(group.gmobs)
	ind := sort.Search(length, func(i int) bool {
		return group.gmobs[i].zDraw >= zDraw
	})

	if ind == 0 {
		group.gmobs = append([]*GameObject{gmob},
			group.gmobs...)
	} else if ind < length {
		group.gmobs = append(group.gmobs[:ind+1],
			group.gmobs[ind:]...)
		group.gmobs[ind] = gmob
	} else {
		group.gmobs = append(group.gmobs, gmob)
	}

	gmob.zDraw = zDraw
}

// addGameObject adds a new game object in the
// draw group computing its position by the priority.
func (group *drawGroup) addGameObject(gmob *GameObject, zDraw float64) error {
	if group.hasGameObject(gmob.Name()) {
		return fmt.Errorf("the game object with name'%s'"+
			" already exists in the draw group", gmob.Name())
	}

	group.insertGameObject(gmob, zDraw)

	return nil
}

// findGameObject returns the game object with the specified name
// or nil if t's not in the group.
func (group *drawGroup) findGameObject(name string) (int, *GameObject) {
	ind := -1
	var gameObject *GameObject

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
	gmobs := []*GameObject{}

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
func (group *drawGroup) addRenderObject(sprite *pixel.Sprite, transform pixel.Matrix, mask color.Color, zDraw float64) {
	rnob := renderObject{
		sprite:    sprite,
		transform: transform,
		mask:      mask,
		zDraw:     zDraw,
	}

	length := len(group.rnobs)
	ind := sort.Search(length, func(i int) bool {
		return group.rnobs[i].zDraw >= zDraw
	})

	if ind == 0 {
		group.rnobs = append([]renderObject{rnob},
			group.rnobs...)
	} else if ind < length {
		group.rnobs = append(group.rnobs[:ind+1],
			group.rnobs[ind:]...)
		group.rnobs[ind] = rnob
	} else {
		group.rnobs = append(group.rnobs, rnob)
	}
}

// addTextObject adds the text object to the draw group
// so it will be drawn in the current frame and removed right after it.
func (group *drawGroup) addTextObject(txt *text.Text, transform pixel.Matrix, mask color.Color, zDraw float64) {
	txob := textObject{
		txt:       txt,
		transform: transform,
		mask:      mask,
		zDraw:     zDraw,
	}

	length := len(group.txobs)
	ind := sort.Search(length, func(i int) bool {
		return group.txobs[i].zDraw >= zDraw
	})

	if ind == 0 {
		group.txobs = append([]textObject{txob},
			group.txobs...)
	} else if ind < length {
		group.txobs = append(group.txobs[:ind+1],
			group.txobs[ind:]...)
		group.txobs[ind] = txob
	} else {
		group.txobs = append(group.txobs, txob)
	}
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

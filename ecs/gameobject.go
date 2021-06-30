package ecs

import (
	"fmt"

	"github.com/alacrity-engine/core/geometry"

	"github.com/faiface/pixel"
)

// GameObject represents a single object
// in the game world which contains components
// to be updated once per frame.
type GameObject struct {
	name       string
	components []Component
	transform  *geometry.Transform
	sprite     *pixel.Sprite
	colorMask  pixel.RGBA
	scene      *Scene
	draw       bool
	destroyed  bool
}

// Destroyed returns true if the game object
// has been destroyed and must not be drawn
// or updated.
func (gmob *GameObject) Destroyed() bool {
	return gmob.destroyed
}

// SetDraw sets if the game object should be drawn.
func (gmob *GameObject) SetDraw(draw bool) {
	gmob.draw = draw
}

// ShouldBeDrawn returns true if the game object
// should be drawn, and false otherwise.
func (gmob *GameObject) ShouldBeDrawn() bool {
	return gmob.draw
}

// Name returns the name of the game object.
func (gmob *GameObject) Name() string {
	return gmob.name
}

// Transform returns the transform to perform
// affine transformations on the game object.
func (gmob *GameObject) Transform() *geometry.Transform {
	return gmob.transform
}

// Start starts all the components
// of the game object.
func (gmob *GameObject) Start() error {
	for _, comp := range gmob.components {
		err := comp.Start()

		if err != nil {
			return err
		}
	}

	gmob.draw = true

	return nil
}

// Update calls update method on all
// game object components.
func (gmob *GameObject) Update() error {
	for _, comp := range gmob.components {
		if comp.Active() {
			err := comp.Update()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Sprite returns the graphical sprite of the
// game object.
func (gmob *GameObject) Sprite() *pixel.Sprite {
	return gmob.sprite
}

// SetSprite sets the sprite for the current game object.
func (gmob *GameObject) SetSprite(sprite *pixel.Sprite) {
	gmob.sprite = sprite
}

// ColorMask returns the value of the color mask
// which is used in drawing sprites.
func (gmob *GameObject) ColorMask() pixel.RGBA {
	return gmob.colorMask
}

// SetColorMask sets of the color mask which is used
// in drawing sprites.
func (gmob *GameObject) SetColorMask(mask pixel.RGBA) {
	gmob.colorMask = mask
}

// Draw the game object onto the target.
func (gmob *GameObject) Draw(target pixel.Target) {
	if gmob.draw && gmob.sprite != nil {
		gmob.sprite.DrawColorMask(target,
			gmob.transform.Data(), gmob.colorMask)
	}
}

// FindComponent searches for the component with the specified name
// and returns it if it exists.
func (gmob *GameObject) FindComponent(name string) (int, Component) {
	ind := -1
	var component Component

	for i, comp := range gmob.components {
		if name == comp.Name() {
			ind = i
			component = comp
			break
		}
	}

	return ind, component
}

// Scene returns the scene where the game object
// resides.
func (gmob *GameObject) Scene() *Scene {
	return gmob.scene
}

// SetScene sets the scene for the game object.
func (gmob *GameObject) SetScene(scene *Scene) {
	gmob.scene = scene
}

// HasComponent returns true if the game object has a
// component with the specified name.
func (gmob *GameObject) HasComponent(name string) bool {
	_, component := gmob.FindComponent(name)

	return component != nil
}

// AddComponent adds the component in the game object.
func (gmob *GameObject) AddComponent(component Component, priority int) error {
	if gmob.HasComponent(component.Name()) {
		return fmt.Errorf("game object '%s'"+
			" already has component '%s'", gmob.Name(), component.Name())
	}

	length := len(gmob.components)

	if length <= 0 || priority >= length {
		gmob.components = append(gmob.components,
			component)
	} else if priority < 0 {
		temp := gmob.components[1:]

		gmob.components = []Component{component}
		gmob.components = append(gmob.components, temp...)
	} else {
		gmob.components = append(gmob.components[:priority+1],
			gmob.components[priority:]...)
		gmob.components[priority] = component
	}

	component.SetGameObject(gmob)

	return nil
}

// RemoveComponent removes the component with the specified
// name from the game object,
func (gmob *GameObject) RemoveComponent(name string) error {
	i, component := gmob.FindComponent(name)

	if component == nil {
		return fmt.Errorf("game object '%s' has no"+
			" component '%s'", gmob.Name(), name)
	}

	err := component.Destroy()

	if err != nil {
		return err
	}

	gmob.components = append(gmob.components[:i],
		gmob.components[i+1:]...)
	component.SetGameObject(nil)

	return nil
}

// ComponentCount returns the number of components
// of the game object.
func (gmob *GameObject) ComponentCount() int {
	return len(gmob.components)
}

// NewGameObject creates a new game object with no components.
func NewGameObject(parent *geometry.Transform, name string, sprite *pixel.Sprite) *GameObject {
	return &GameObject{
		name:       name,
		components: []Component{},
		transform:  geometry.NewTransform(parent, pixel.IM),
		sprite:     sprite,
		draw:       false,
		colorMask:  pixel.Alpha(1.0),
	}
}

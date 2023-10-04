package engine

import (
	"fmt"

	cmath "github.com/alacrity-engine/core/math"
	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
)

// TODO: add an opportunity for a
// gameobject to inherit a sprite
// from another game object by link.

// GameObject represents a single object
// in the game world which contains components
// to be updated once per frame.
type GameObject struct {
	name       string
	components map[string]Component
	transform  *geometry.Transform
	sprite     *render.Sprite
	scene      *Scene
	draw       bool
	destroyed  bool
	zUpdate    cmath.Fixed
	started    bool
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
	if gmob.started {
		return nil
	}

	for _, comp := range gmob.components {
		err := comp.Start()

		if err != nil {
			return err
		}
	}

	gmob.draw = true
	gmob.started = true

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
func (gmob *GameObject) Sprite() *render.Sprite {
	return gmob.sprite
}

// SetSprite sets the sprite for the current game object.
func (gmob *GameObject) SetSprite(sprite *render.Sprite) {
	gmob.sprite = sprite
}

// Draw the game object onto the target.
func (gmob *GameObject) Draw() error {
	if gmob.draw && gmob.sprite != nil {
		return gmob.sprite.Draw(gmob.transform)
	}

	return nil
}

// FindComponent searches for the component with the specified name
// and returns it if it exists.
func (gmob *GameObject) FindComponent(name string) Component {
	return gmob.components[name]
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
	component := gmob.FindComponent(name)

	return component != nil
}

// AddComponent adds the component in the game object.
func (gmob *GameObject) AddComponent(component Component, priority int) error {
	regComp, ok := component.(RegisteredComponent)

	if !ok {
		return fmt.Errorf("the component can't be registered")
	}

	typeID := regComp.TypeID()

	if _, ok := gmob.components[typeID]; ok {
		return fmt.Errorf(
			"the game object '%s' already has a '%s' component",
			gmob.name, typeID)
	}

	gmob.components[typeID] = component
	component.SetGameObject(gmob)

	return nil
}

// RemoveComponent removes the component with the specified
// name from the game object,
func (gmob *GameObject) RemoveComponent(component Component) error {
	regComp, ok := component.(RegisteredComponent)

	if !ok {
		return fmt.Errorf("the component can't be registered")
	}

	if !gmob.HasComponent(regComp.TypeID()) {
		return RaiseErrorNoComponentOnGameObject(gmob, regComp.TypeID())
	}

	err := component.Destroy()

	if err != nil {
		return err
	}

	delete(gmob.components, regComp.TypeID())
	component.SetGameObject(nil)

	return nil
}

// ComponentCount returns the number of components
// of the game object.
func (gmob *GameObject) ComponentCount() int {
	return len(gmob.components)
}

// NewGameObject creates a new game object with no components.
func NewGameObject(parent *geometry.Transform, name string, sprite *render.Sprite) *GameObject {
	return &GameObject{
		name:       name,
		components: map[string]Component{},
		transform:  geometry.NewTransform(parent),
		sprite:     sprite,
		draw:       false,
	}
}

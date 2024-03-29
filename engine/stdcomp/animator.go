package stdcomp

import (
	"fmt"

	"github.com/alacrity-engine/core/anim"
	"github.com/alacrity-engine/core/engine"
)

// Dummy is used to
// stop the current animation
// but don't play any else.
const Dummy = ""

// Animator is a component which controls
// what animation to play.
type Animator struct {
	engine.BaseComponent
	currentAnimation string
	animations       map[string]*anim.Animation
}

// CurrentAnimation returns the animation currently
// being played.
func (animator *Animator) CurrentAnimation() string {
	return animator.currentAnimation
}

// PlayAnimation plays the specified animation.
func (animator *Animator) PlayAnimation(anim string) error {
	if _, ok := animator.animations[anim]; anim != Dummy && !ok {
		return fmt.Errorf("animation with name '%s' doesn't exist",
			anim)
	}

	animator.StopAnimation()

	if anim != "" {
		animator.animations[anim].SetSprite(animator.GameObject().Sprite())
		animator.animations[anim].Start()
	}

	animator.currentAnimation = anim

	return nil
}

// StopAnimation stops the current animation.
func (animator *Animator) StopAnimation() {
	if animator.currentAnimation != Dummy {
		animator.animations[animator.currentAnimation].Stop()
	}
}

// AnimationActive returns true if the animation with the given name
// is currently being played.
func (animator *Animator) AnimationActive(anim string) (bool, error) {
	if _, ok := animator.animations[anim]; !ok {
		return false, fmt.Errorf("animation with name '%s' doesn't exist",
			anim)
	}

	return animator.animations[anim].Active(), nil
}

// Start starts the animator component.
func (animator *Animator) Start() error {
	return nil
}

// Update sets the current animation sprite to the game object.
func (animator *Animator) Update() error {
	return animator.animations[animator.currentAnimation].Update()
}

// Destroy destroys the component and
// stops the current animation.
func (animator *Animator) Destroy() error {
	animator.StopAnimation()

	for _, anim := range animator.animations {
		err := anim.Dispose()

		if err != nil {
			return err
		}
	}

	return nil
}

// AddAnimation adds the animation with the given name in the animator.
func (animator *Animator) AddAnimation(name string, anim *anim.Animation) error {
	if _, ok := animator.animations[name]; ok {
		return fmt.Errorf("animation with name '%s' already exists",
			name)
	}

	animator.animations[name] = anim

	return nil
}

// RemoveAnimation removes the animation from the animator by
// its name.
func (animator *Animator) RemoveAnimation(name string) error {
	if _, ok := animator.animations[name]; !ok {
		return fmt.Errorf("animation with name '%s' doesn't exist",
			name)
	}

	err := animator.PlayAnimation(Dummy)

	if err != nil {
		return err
	}

	delete(animator.animations, name)

	return nil
}

// NewAnimator creates a new animator with
// the given name.
func NewAnimator(name string) *Animator {
	animator := &Animator{
		currentAnimation: Dummy,
		animations:       map[string]*anim.Animation{},
	}

	return animator
}

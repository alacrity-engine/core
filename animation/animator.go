package animation

import (
	"fmt"

	"github.com/alacrity-engine/core/ecs"
)

// Dummy is used to
// stop the current animation
// but don't play any else.
const Dummy = ""

// Animator is a component which controls
// what animation to play.
type Animator struct {
	ecs.BaseComponent
	currentAnimation string
	animations       map[string]*Animation
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
	currentSprite := animator.
		animations[animator.currentAnimation].GetCurrentSprite()
	animator.GameObject().SetSprite(currentSprite)

	return nil
}

// Destroy destroys the component and
// stops the current animation.
func (animator *Animator) Destroy() error {
	animator.StopAnimation()

	return nil
}

// AddAnimation adds the animation with the given name in the animator.
func (animator *Animator) AddAnimation(name string, anim *Animation) error {
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

	animator.PlayAnimation(Dummy)
	delete(animator.animations, name)

	return nil
}

// NewAnimator creates a new animator with
// the given name.
func NewAnimator(name string) *Animator {
	animator := &Animator{
		currentAnimation: Dummy,
		animations:       map[string]*Animation{},
	}

	animator.SetName(name)

	return animator
}

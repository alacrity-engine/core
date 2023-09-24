package anim

import (
	"fmt"
	"time"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
)

// Animation represent a single
// animation made of sprites.
type Animation struct {
	frames        []geometry.Rect
	delays        []time.Duration
	cancel        chan struct{}
	currentFrame  int
	texture       *render.Texture
	currentSprite *render.Sprite
	active        bool
	loop          bool
}

// Loop returns true if the animation
// is looped, and false otherwise.
func (anim *Animation) Loop() bool {
	return anim.loop
}

// SetLoop sets the animation to be looped
// or not.
func (anim *Animation) SetLoop(loop bool) {
	anim.loop = loop
}

// Active returns true if the animation
// is currently being played, and false
// otherwise.
func (anim *Animation) Active() bool {
	return anim.active
}

func (anim *Animation) SetSprite(sprite *render.Sprite) error {
	if anim.texture != sprite.Texture() {
		return fmt.Errorf(
			"the sprite must have the same texture as the animtion")
	}

	anim.currentSprite = sprite

	return nil
}

// Start starts playing animation.
func (anim *Animation) Start() {
	if anim.active {
		return
	}

	anim.currentFrame = 0
	anim.setFrame(anim.currentFrame)
	anim.active = true

	timeout := time.After(anim.delays[anim.currentFrame])
	anim.cancel = make(chan struct{})

	go anim.process(timeout, anim.cancel)
}

func (anim *Animation) process(timeout <-chan time.Time, cancel <-chan struct{}) {
	for {
		select {
		case <-timeout:
			if anim.currentFrame+1 >= len(anim.frames) {
				if !anim.loop {
					return
				}

				anim.currentFrame = 0
			} else {
				anim.currentFrame++
			}

			timeout = time.After(anim.delays[anim.currentFrame])

		case <-cancel:
			return
		}
	}
}

func (anim *Animation) setFrame(ind int) error {
	if anim.currentSprite == nil {
		return fmt.Errorf(
			"no sprite specified for the animation")
	}

	err := anim.currentSprite.SetTargetArea(
		anim.frames[ind])

	if err != nil {
		return err
	}

	return nil
}

func (anim *Animation) Update() error {
	return anim.setFrame(anim.currentFrame)
}

// Stop stops playing animation.
func (anim *Animation) Stop() {
	if anim.active {
		anim.active = false
		anim.cancel <- struct{}{}
		anim.cancel = nil
	}
}

func (anim *Animation) Dispose() error {
	anim.Stop()

	return nil
}

// GetCurrentSprite returns a new sprite for the animation frame
// played at the moment.
func (anim *Animation) GetCurrentSprite() *render.Sprite {
	return anim.currentSprite
}

// NewAnimation creates a new animation
// out of frames and their delays.
func NewAnimation(
	texture *render.Texture,
	frames []geometry.Rect,
	delays []time.Duration,
	loop bool,
) (*Animation, error) {
	anim := &Animation{
		frames:       make([]geometry.Rect, len(frames)),
		delays:       make([]time.Duration, len(delays)),
		currentFrame: 0,
		texture:      texture,
		active:       false,
		loop:         loop,
	}

	copy(anim.frames, frames)
	copy(anim.delays, delays)

	return anim, nil
}

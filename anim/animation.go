package anim

import (
	"fmt"
	"time"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
)

// TODO: rewrite it without
// goroutines and channels.

// Animation represent a single
// animation made of sprites.
type Animation struct {
	frames        []geometry.Rect
	delays        []time.Duration
	timeout       <-chan time.Time
	cancel        chan interface{}
	errCh         chan error
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

func (anim *Animation) LastError() error {
	select {
	case err := <-anim.errCh:
		return err

	default:
		return nil
	}
}

// Start starts playing animation.
func (anim *Animation) Start() {
	if anim.active {
		return
	}

	anim.currentFrame = 0
	anim.setFrame(anim.currentFrame)
	anim.timeout = time.After(anim.delays[anim.currentFrame])
	anim.cancel = make(chan interface{})
	anim.active = true

	go anim.process()
}

func (anim *Animation) process() {
	for {
		select {
		case <-anim.timeout:
			anim.currentFrame++

			if anim.currentFrame >= len(anim.frames) {
				if !anim.loop {
					anim.stopAnimating()

					return
				}

				anim.currentFrame = 0
			}

			err := anim.setFrame(anim.currentFrame)

			if err != nil {
				select {
				case <-anim.errCh:
					anim.errCh <- err

				default:
					anim.errCh <- err
				}
			}

			anim.timeout = time.After(anim.delays[anim.currentFrame])

		case <-anim.cancel:
			anim.stopAnimating()

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

func (anim *Animation) stopAnimating() {
	anim.active = false

	close(anim.cancel)
	anim.cancel = nil
	anim.timeout = nil
}

// Stop stops playing animation.
func (anim *Animation) Stop() {
	if anim.active {
		go func() {
			anim.cancel <- true
		}()
	}
}

func (anim *Animation) Dispose() error {
	anim.Stop()
	close(anim.errCh)

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
		errCh:        make(chan error, 1),
		currentFrame: 0,
		texture:      texture,
		active:       false,
		loop:         loop,
	}

	copy(anim.frames, frames)
	copy(anim.delays, delays)

	return anim, nil
}

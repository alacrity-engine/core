package anim

import (
	"time"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
	codec "github.com/alacrity-engine/resource-codec"
)

// Animation represent a single
// animation made of sprites.
type Animation struct {
	frames        []geometry.Rect
	delays        []time.Duration
	timeout       <-chan time.Time
	cancel        chan interface{}
	currentFrame  int
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

// Start starts playing animation.
func (anim *Animation) Start() {
	if anim.active {
		return
	}

	anim.currentFrame = 0
	anim.setSprite(anim.currentFrame)
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

			anim.setSprite(anim.currentFrame)
			anim.timeout = time.After(anim.delays[anim.currentFrame])

		case <-anim.cancel:
			anim.stopAnimating()

			return
		}
	}
}

func (anim *Animation) setSprite(ind int) {
	anim.currentSprite.SetTargetArea(
		anim.frames[ind])
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

// GetCurrentSprite returns a new sprite for the animation frame
// played at the moment.
func (anim *Animation) GetCurrentSprite() *render.Sprite {
	return anim.currentSprite
}

// NewAnimation creates a new animation
// out of frames and their delays.
func NewAnimation(
	spritesheet *render.Texture,
	shaderProgram *render.ShaderProgram,
	vertexDrawMode,
	colorDrawMode render.DrawMode,
	frames []geometry.Rect,
	delays []time.Duration,
	loop bool,
) (*Animation, error) {
	animSprite, err := render.NewSpriteFromTextureAndProgram(
		vertexDrawMode, render.DrawModeDynamic,
		colorDrawMode, spritesheet,
		shaderProgram, geometry.Rect{})

	if err != nil {
		return nil, err
	}

	anim := &Animation{
		frames:        make([]geometry.Rect, len(frames)),
		delays:        make([]time.Duration, len(delays)),
		currentFrame:  0,
		currentSprite: animSprite,
		active:        false,
		loop:          loop,
	}

	copy(anim.frames, frames)
	copy(anim.delays, delays)

	return anim, nil
}

// NewAnimationWithExistingSprite cerates
// a new animation with the existing sprite.
func NewAnimationWithExistingSprite(
	sprite *render.Sprite,
	frames []geometry.Rect,
	delays []time.Duration,
	loop bool,
) *Animation {
	anim := &Animation{
		frames:        make([]geometry.Rect, len(frames)),
		delays:        make([]time.Duration, len(delays)),
		currentFrame:  0,
		currentSprite: sprite,
		active:        false,
		loop:          loop,
	}

	copy(anim.frames, frames)
	copy(anim.delays, delays)

	return anim
}

func NewAnimationFromPictureAndData(
	picture *codec.Picture,
	filter render.TextureFiltering,
	vertexDrawMode,
	colorDrawMode render.DrawMode,
	shaderProgram *render.ShaderProgram,
	frames []geometry.Rect,
	delays []time.Duration,
	loop bool,
) (*Animation, error) {
	texture := render.NewTextureFromPicture(picture, filter)
	animSprite, err := render.NewSpriteFromTextureAndProgram(
		vertexDrawMode, render.DrawModeDynamic,
		colorDrawMode, texture,
		shaderProgram, geometry.Rect{})

	if err != nil {
		return nil, err
	}

	anim := &Animation{
		frames:        make([]geometry.Rect, len(frames)),
		delays:        make([]time.Duration, len(delays)),
		currentFrame:  0,
		currentSprite: animSprite,
		active:        false,
		loop:          loop,
	}

	copy(anim.frames, frames)
	copy(anim.delays, delays)

	return anim, nil
}

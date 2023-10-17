package resources

import (
	"fmt"

	"github.com/alacrity-engine/core/render"
	codec "github.com/alacrity-engine/resource-codec"
	"github.com/golang/freetype/truetype"
)

// resourceBuffer stores all the resources
// ever loaded by the game scene. The resource
// loader will at first try to look up the resource
// in the buffer and if it isn't loaded the loader
// will take the resource from the resource file.
type resourceBuffer struct {
	pictures       map[string]*render.Picture
	animations     map[string]*codec.AnimationData
	fonts          map[string]*truetype.Font
	audio          map[string][]byte
	textures       map[string]*render.Texture
	shaders        map[string]*render.Shader
	shaderPrograms map[string]*render.ShaderProgram
	spritesheets   map[string]*codec.SpritesheetData
}

// putPicture puts the picture in the buffer.
func (rb *resourceBuffer) putPicture(name string, pic *render.Picture) error {
	if _, ok := rb.pictures[name]; ok {
		return RaiseErrorPictureAlreadyExists(name)
	}

	rb.pictures[name] = pic

	return nil
}

// takePicture takes the picture from the buffer.
func (rb *resourceBuffer) takePicture(name string) (*render.Picture, error) {
	if _, ok := rb.pictures[name]; !ok {
		return nil, RaiseErrorPictureDoesntExist(name)
	}

	return rb.pictures[name], nil
}

func (rb *resourceBuffer) putSpritesheet(name string, ss *codec.SpritesheetData) error {
	if _, ok := rb.spritesheets[name]; ok {
		return fmt.Errorf("the '%s' spritesheet already exists", name)
	}

	rb.spritesheets[name] = ss

	return nil
}

func (rb *resourceBuffer) takeSpritesheet(name string) (*codec.SpritesheetData, error) {
	if _, ok := rb.spritesheets[name]; !ok {
		return nil, RaiseErrorSpritesheetDoesntExist(name)
	}

	return rb.spritesheets[name], nil
}

func (rb *resourceBuffer) putTexture(name string, texture *render.Texture) error {
	if _, ok := rb.textures[name]; ok {
		return fmt.Errorf("the '%s' texture already exists", name)
	}

	rb.textures[name] = texture

	return nil
}

func (rb *resourceBuffer) takeTexture(name string) (*render.Texture, error) {
	if _, ok := rb.textures[name]; !ok {
		return nil, RaiseErrorTextureDoesntExist(name)
	}

	return rb.textures[name], nil
}

// putAnimation puts the animation in the buffer.
func (rb *resourceBuffer) putAnimation(name string, anim *codec.AnimationData) error {
	if _, ok := rb.animations[name]; ok {
		return RaiseErrorAnimationAlreadyExists(name)
	}

	rb.animations[name] = anim

	return nil
}

// takeAnimation takes the animation from the buffer.
func (rb *resourceBuffer) takeAnimation(name string) (*codec.AnimationData, error) {
	if _, ok := rb.animations[name]; !ok {
		return nil, RaiseErrorAnimationDoesntExist(name)
	}

	return rb.animations[name], nil
}

// putFont puts the font in the buffer.
func (rb *resourceBuffer) putFont(name string, fnt *truetype.Font) error {
	if _, ok := rb.fonts[name]; ok {
		return RaiseErrorFontAlreadyExists(name)
	}

	rb.fonts[name] = fnt

	return nil
}

// takeFont takes the font from the buffer.
func (rb *resourceBuffer) takeFont(name string) (*truetype.Font, error) {
	if _, ok := rb.fonts[name]; !ok {
		return nil, RaiseErrorFontDoesntExist(name)
	}

	return rb.fonts[name], nil
}

// putAudio puts the audio in the buffer.
func (rb *resourceBuffer) putAudio(name string, audio []byte) error {
	if _, ok := rb.audio[name]; ok {
		return RaiseErrorAudioAlreadyExists(name)
	}

	rb.audio[name] = audio

	return nil
}

// takeAudio takes the audio from the buffer.
func (rb *resourceBuffer) takeAudio(name string) ([]byte, error) {
	if _, ok := rb.audio[name]; !ok {
		return nil, RaiseErrorAudioDoesntExist(name)
	}

	return rb.audio[name], nil
}

// newResourceBuffer creates a new resource buffer
// to store every resource ever loaded by the loader.
func newResourceBuffer() *resourceBuffer {
	return &resourceBuffer{
		pictures:   map[string]*render.Picture{},
		animations: map[string]*codec.AnimationData{},
		fonts:      map[string]*truetype.Font{},
		audio:      map[string][]byte{},
	}
}

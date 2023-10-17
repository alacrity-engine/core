package typography

import (
	"fmt"

	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
)

// TODO: automatically create a batch
// for each loaded font atlas.

type Atlas struct {
	Frames      map[rune]geometry.Rect
	CharTexture *render.Texture
}

func (atlas *Atlas) CreateCharSprite(
	vertexDrawMode, textureDrawMode, colorDrawMode render.DrawMode,
	texture *render.Texture, shaderProgram *render.ShaderProgram,
	char rune,
) (*render.Sprite, error) {
	if atlas.CharTexture == nil {
		return nil, fmt.Errorf("the character texture is nil")
	}

	frame, ok := atlas.Frames[char]

	if !ok {
		return nil, fmt.Errorf("no character entry for '%c'", char)
	}

	sprite, err := render.NewSpriteFromTextureAndProgram(
		vertexDrawMode, textureDrawMode, colorDrawMode,
		atlas.CharTexture, shaderProgram, frame)

	if err != nil {
		return nil, err
	}

	return sprite, nil
}

func (atlas *Atlas) SetCharForSprite(sprite *render.Sprite, char rune) error {
	if sprite.Texture() != atlas.CharTexture {
		return fmt.Errorf("textures don't match")
	}

	frame, ok := atlas.Frames[char]

	if !ok {
		return fmt.Errorf("no character entry for '%c'", char)
	}

	err := sprite.SetTargetArea(frame)

	if err != nil {
		return err
	}

	return nil
}

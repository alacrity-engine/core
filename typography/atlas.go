package typography

import (
	"fmt"

	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
)

type Atlas struct {
	Glyphs      map[rune]Glyph
	CharTexture *render.Texture
}

type Glyph struct {
	Dot     geometry.Vec
	Frame   geometry.Rect
	Advance float64
}

func (atlas *Atlas) CreateCharSprite(
	vertexDrawMode, textureDrawMode, colorDrawMode render.DrawMode,
	texture *render.Texture, shaderProgram *render.ShaderProgram,
	char rune,
) (*render.Sprite, error) {
	if atlas.CharTexture == nil {
		return nil, fmt.Errorf("the character texture is nil")
	}

	glyph, ok := atlas.Glyphs[char]

	if !ok {
		return nil, fmt.Errorf("no character entry for '%c'", char)
	}

	sprite, err := render.NewSpriteFromTextureAndProgram(
		vertexDrawMode, textureDrawMode, colorDrawMode,
		atlas.CharTexture, shaderProgram, glyph.Frame)

	if err != nil {
		return nil, err
	}

	return sprite, nil
}

func (atlas *Atlas) SetCharForSprite(sprite *render.Sprite, char rune) error {
	if sprite.Texture() != atlas.CharTexture {
		return fmt.Errorf("textures don't match")
	}

	glyph, ok := atlas.Glyphs[char]

	if !ok {
		return fmt.Errorf("no character entry for '%c'", char)
	}

	err := sprite.SetTargetArea(glyph.Frame)

	if err != nil {
		return err
	}

	return nil
}

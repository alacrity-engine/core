package stdcomp

import (
	"github.com/alacrity-engine/core/engine"
	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/typography"
)

// TODO: add all the letter
// sprites to the font batch.

type Caption struct {
	engine.BaseComponent
	text        []rune
	caret       geometry.Vec
	charSprites []*render.Sprite
	batch       *render.Batch     `iris:"exported"`
	atlas       *typography.Atlas `iris:"exported"`
	width       int               `iris:"exported"`
	height      int               `iris:"exported"`
}

func (caption *Caption) Start() error {
	halfDiag := geometry.V(
		float64(caption.width),
		float64(caption.height),
	)
	caption.caret = caption.
		GameObject().Transform().Position().
		Add(halfDiag.Scaled(-1))

	err := caption.GameObject().DelegateDrawing(caption)

	if err != nil {
		return err
	}

	return nil
}

func (caption *Caption) Update() error {
	return nil
}

func (caption *Caption) Draw(transform *geometry.Transform) error {
	return nil
}

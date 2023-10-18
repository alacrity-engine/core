package stdcomp

import (
	"fmt"
	"unsafe"

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
	batch       *render.Batch
	atlas       *typography.Atlas `iris:"exported"`
	width       int               `iris:"exported"`
	height      int               `iris:"exported"`
}

func (caption *Caption) Start() error {
	halfDiag := geometry.V(float64(caption.width), float64(caption.height))
	caption.caret = caption.GameObject().Transform().Position().Add(halfDiag.Scaled(-1))

	atlasAddress := uintptr(unsafe.Pointer(caption.atlas))
	batchID := fmt.Sprintf("__tab%X", atlasAddress)
	var err error
	caption.batch, err = engine.BatchByName(batchID)

	if err != nil {
		return err
	}

	return nil
}

func (caption *Caption) Update() error {
	return nil
}

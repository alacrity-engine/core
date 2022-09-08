package typography

import (
	"fmt"

	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

// AtlasBuffer is a storage for
// all the created text atlases.
type AtlasBuffer struct {
	atlases map[string]*text.Atlas
}

// Atlas returns the atlas by name.
func (ab *AtlasBuffer) Atlas(name string) (*text.Atlas, error) {
	if _, ok := ab.atlases[name]; !ok {
		return nil, fmt.Errorf("no atlas named '%s'", name)
	}

	return ab.atlases[name], nil
}

// CreateAtlas creates a new atlas
// wirth specified name.
func (ab *AtlasBuffer) CreateAtlas(name string, face font.Face, runeSets ...[]rune) error {
	if _, ok := ab.atlases[name]; ok {
		return fmt.Errorf("atlas named '%s' already exists", name)
	}

	atlas := text.NewAtlas(face, runeSets...)
	ab.atlases[name] = atlas

	return nil
}

// NewAtlasBuffer creates a new
// storage for text atlases.
func NewAtlasBuffer() *AtlasBuffer {
	return &AtlasBuffer{
		atlases: map[string]*text.Atlas{},
	}
}

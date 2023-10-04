package definitions

import (
	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
)

type Prefab struct {
	Name          string
	TransformRoot *TransformDefinition
}

type TransformDefinition struct {
	Position geometry.Vec
	Angle    float64
	Scale    geometry.Vec
	Gmob     *GameObjectDefinition
	Children []*TransformDefinition
}

type GameObjectDefinition struct {
	Name       string
	ZUpdate    float64
	Components []*ComponentDefinition
	Sprite     *SpriteDefinition
	Draw       bool
}

type ComponentDefinition struct {
	TypeName string
	Active   bool
	Data     map[string]interface{}
}

type SpriteDefinition struct {
	ColorMask       render.ColorMask
	TargetArea      geometry.Rect
	VertexDrawMode  render.DrawMode
	TextureDrawMode render.DrawMode
	ColorDrawMode   render.DrawMode
	ShaderProgramID string
	TextureID       string
	CanvasID        string
	BatchID         string
}

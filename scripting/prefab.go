package scripting

import (
	"github.com/alacrity-engine/core/geometry"
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
	DrawMode        render.DrawMode
	ShaderProgramID string
	TextureID       string
	CanvasID        string
	BatchID         string
}

type GameObjectPointer struct {
	Name string
}

type ComponentPointer struct {
	GmobName string
	CompType string
}

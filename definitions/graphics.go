package definitions

type BatchDefinition struct {
	Name      string
	CanvasID  string
	TextureID string
	ZMin      float32
	ZMax      float32
}

type CanvasDefinition struct {
	Name  string
	DrawZ int
}

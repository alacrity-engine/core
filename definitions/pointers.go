package definitions

const (
	ResourceTypeAnimation   = "animation"
	ResourceTypeAudio       = "audio"
	ResourceTypePicture     = "picture"
	ResourceTypeTexture     = "texture"
	ResourceTypeFont        = "font"
	ResourceTypeSpritesheet = "spritesheet"
)

type GameObjectPointer struct {
	Name string
}

type ComponentPointer struct {
	GmobName string
	CompType string
}

type ResourcePointer struct {
	ResourceType string
	ResourceID   string
}

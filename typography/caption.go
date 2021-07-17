package typography

import (
	"fmt"
	"image/color"

	"github.com/alacrity-engine/core/ecs"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

// Caption is a text
// attached to the game object.
type Caption struct {
	ecs.BaseComponent
	orig      pixel.Vec
	atlas     *text.Atlas
	str       string
	txt       *text.Text
	target    string
	colorMask color.Color
	zDraw     float64
}

// String returns the contents of the caption.
func (c *Caption) String() string {
	return c.str
}

// SetString sets the contents of the caption.
func (c *Caption) SetString(str string) {
	c.str = str
}

// ColorMask returns the color of the text.
func (c *Caption) ColorMask() color.Color {
	return c.colorMask
}

// SetColorMask sets the color of the text.
func (c *Caption) SetColorMask(mask color.Color) {
	c.colorMask = mask
}

// Atlas returns the atlas
// of the text object.
func (c *Caption) Atlas() *text.Atlas {
	return c.atlas
}

// SetAtlas sets the atlas
// for the text object.
func (c *Caption) SetAtlas(atlas *text.Atlas) {
	c.txt = text.New(c.orig, atlas)
	c.atlas = atlas
}

// Origin returns the origin
// of the text object.
func (c *Caption) Origin() pixel.Vec {
	return c.orig
}

// SetOrigin sets the origin
// for the text object.
func (c *Caption) SetOrigin(orig pixel.Vec) {
	c.txt = text.New(orig, c.atlas)
	c.orig = orig
}

// Start does nothing.
func (c *Caption) Start() error {
	return nil
}

// Update fills the text object with
// the content and sends it to be rendered.
func (c *Caption) Update() error {
	layout := c.GameObject().Scene().DrawLayout()
	fmt.Fprint(c.txt, c.str)

	err := layout.RenderTextOnTarget(
		c.target, c.txt, c.GameObject().
			Transform().Data(), c.colorMask, c.zDraw)

	if err != nil {
		return err
	}

	return nil
}

// Destroy does nothing.
func (c *Caption) Destroy() error {
	return nil
}

// NewCaption returns a new text object to
// be rendered on the game scene.
func NewCaption(orig pixel.Vec, name, target string, atlas *text.Atlas, zDraw float64) *Caption {
	caption := &Caption{
		target: target,
		txt:    text.New(orig, atlas),
		orig:   orig,
		zDraw:  zDraw,
	}

	caption.SetName(name)

	return caption
}

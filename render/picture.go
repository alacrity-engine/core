package render

import "github.com/alacrity-engine/core/geometry"

type Picture struct {
	Width  int32
	Height int32
	Pix    []byte
}

// GetSpritesheetFrames returns the set of rectangles
// corresponding to the frames of the spritesheet.
func (spritesheet *Picture) GetSpritesheetFrames(width, height int) []geometry.Rect {
	frames := make([]geometry.Rect, 0)
	pixelWidth := float64(spritesheet.Width)
	pixelHeight := float64(spritesheet.Height)
	dw := pixelWidth / float64(width)
	dh := pixelHeight / float64(height)

	for y := pixelHeight; y > 0; y -= dh {
		for x := 0.0; x < pixelWidth; x += dw {
			frame := geometry.R(x, y-dh, x+dw, y)
			frames = append(frames, frame)
		}
	}

	return frames
}

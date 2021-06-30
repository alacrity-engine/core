package system

import (
	"github.com/faiface/pixel/pixelgl"
)

// ButtonPressed returns true if the button is currently pressed down.
func ButtonPressed(button pixelgl.Button) bool {
	return win.Pressed(button)
}

// ButtonJustPressed returns true if the button has just been pressed down.
func ButtonJustPressed(button pixelgl.Button) bool {
	return win.JustPressed(button)
}

// ButtonJustReleased returns true if the button has just been released.
func ButtonJustReleased(button pixelgl.Button) bool {
	return win.JustReleased(button)
}

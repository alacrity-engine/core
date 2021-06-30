package system

import (
	"github.com/faiface/pixel/pixelgl"
)

// Win is the current window
// of the game.
var win *pixelgl.Window

// Window returns the current
// window of the game.
func Window() *pixelgl.Window {
	return win
}

// SetWindow sets the window to read inputs from.
func SetWindow(window *pixelgl.Window) {
	win = window
}

// Resolution returns width and height
// of the current window.
func Resolution() (int, int) {
	return int(win.Bounds().W()), int(win.Bounds().H())
}

// VSyncEnabled detects if vertical
// synchronization is enabled.
func VSyncEnabled() bool {
	return win.VSync()
}

// SetTitle sets a new title for
// the window.
func SetTitle(title string) {
	win.SetTitle(title)
}

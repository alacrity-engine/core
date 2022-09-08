package system

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type window struct {
	winHandler *glfw.Window
}

var (
	mainWindow *window
)

func (win *window) buttonPressed(button Button) bool {
	return win.winHandler.GetKey(glfw.Key(button)) == glfw.Press
}

func InitializeWindow(title string, width, height int, fullscreen, vsync bool) error {
	err := glfw.Init()

	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var monitor *glfw.Monitor

	if fullscreen {
		monitors := glfw.GetMonitors()

		if len(monitors) > 0 {
			monitor = monitors[0]
			SetDisplayResolution(width, height)
		} else {
			return fmt.Errorf("no monitors")
		}
	}

	win, err := glfw.CreateWindow(width, height, title, monitor, nil)

	if err != nil {
		return err
	}

	win.MakeContextCurrent()
	win.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	win.SetInputMode(glfw.StickyKeysMode, glfw.True)

	if !vsync {
		glfw.SwapInterval(0)
	}

	mainWindow = &window{
		winHandler: win,
	}

	return nil
}

// ButtonPressed returns true if the button is currently pressed down.
func ButtonPressed(button Button) bool {
	return mainWindow.buttonPressed(button)
}

func ShouldClose() bool {
	return mainWindow.winHandler.ShouldClose()
}

func TickLoop() {
	mainWindow.winHandler.SwapBuffers()
	glfw.PollEvents()
}

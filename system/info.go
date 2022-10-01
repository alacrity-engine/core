package system

import "github.com/go-gl/gl/v4.6-core/gl"

// Renderer returns the name of the renderer.
func Renderer() string {
	return gl.GoStr(gl.GetString(gl.RENDERER))
}

// Vendor returns the name of the renderer vendor.
func Vendor() string {
	return gl.GoStr(gl.GetString(gl.VENDOR))
}

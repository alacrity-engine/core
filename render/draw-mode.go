package render

import "github.com/go-gl/gl/v4.6-core/gl"

type DrawMode uint32

const (
	DrawModeStatic  DrawMode = gl.STATIC_DRAW
	DrawModeDynamic DrawMode = gl.DYNAMIC_DRAW
	DrawModeStream  DrawMode = gl.STREAM_DRAW
)

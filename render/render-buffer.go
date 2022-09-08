package render

import "github.com/go-gl/gl/v4.6-core/gl"

type RenderBufferType uint32

const (
	RenderBufferTypeDepth            RenderBufferType = gl.DEPTH_COMPONENT
	RenderBufferTypeDepth24Stencil8  RenderBufferType = gl.DEPTH24_STENCIL8
	RenderBufferTypeDepth32fStencil8 RenderBufferType = gl.DEPTH32F_STENCIL8
)

type RenderBuffer struct {
	glHandler uint32
	typ       RenderBufferType
}

func NewRenderBuffer(typ RenderBufferType, width, height int) *RenderBuffer {
	var handler uint32
	gl.GenRenderbuffers(1, &handler)
	gl.BindRenderbuffer(gl.RENDERBUFFER, handler)
	gl.RenderbufferStorage(gl.RENDERBUFFER, uint32(typ),
		int32(width), int32(height))
	gl.BindBuffer(gl.RENDERBUFFER, 0)

	return &RenderBuffer{
		glHandler: handler,
		typ:       typ,
	}
}

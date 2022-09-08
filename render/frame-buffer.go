package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type FrameBufferAttachmentType uint32

const (
	FrameBufferAttachmentTypeDepth        FrameBufferAttachmentType = gl.DEPTH_ATTACHMENT
	FrameBufferAttachmentTypeDepthStencil FrameBufferAttachmentType = gl.DEPTH_STENCIL_ATTACHMENT
)

type FrameBuffer struct {
	glHandler    uint32
	texture      *Texture
	renderBuffer *RenderBuffer
}

func NewFrameBufferWithTextureAndRenderBuffer(texture *Texture, renderBuffer *RenderBuffer) (*FrameBuffer, error) {
	if texture == nil || texture.glHandler == 0 {
		return nil, fmt.Errorf("no texture supplied")
	}

	if renderBuffer == nil || renderBuffer.glHandler == 0 {
		return nil, fmt.Errorf("no render buffer supplied")
	}

	var attachmentType FrameBufferAttachmentType

	switch renderBuffer.typ {
	case RenderBufferTypeDepth:
		attachmentType = FrameBufferAttachmentTypeDepth

	case RenderBufferTypeDepth24Stencil8, RenderBufferTypeDepth32fStencil8:
		attachmentType = FrameBufferAttachmentTypeDepthStencil

	default:
		return nil, fmt.Errorf(
			"invalid render buffer type: %d", renderBuffer.typ)
	}

	var handler uint32
	gl.GenFramebuffers(1, &handler)
	gl.BindFramebuffer(gl.FRAMEBUFFER, handler)
	defer gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	gl.BindTexture(gl.TEXTURE_2D, texture.glHandler)
	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, texture.glHandler, 0)
	defer gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.BindRenderbuffer(gl.RENDERBUFFER, renderBuffer.glHandler)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, uint32(attachmentType), gl.RENDERBUFFER, renderBuffer.glHandler)
	defer gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		gl.DeleteFramebuffers(1, &handler)

		return nil, fmt.Errorf(
			"frame buffer creation error: %d", status)
	}

	return &FrameBuffer{
		glHandler:    handler,
		texture:      texture,
		renderBuffer: renderBuffer,
	}, nil
}

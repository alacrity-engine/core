package render

import "github.com/go-gl/gl/v4.6-core/gl"

type SpriteTextureSlot uint32

const (
	SpriteTextureSlotMainTexture SpriteTextureSlot = gl.TEXTURE0
)

type BatchTextureSlot uint32

const (
	BatchTextureSlotMainTexture          BatchTextureSlot = gl.TEXTURE0
	BatchTextureSlotModelsBuffer         BatchTextureSlot = gl.TEXTURE1
	BatchTextureSlotShouldDrawBuffer     BatchTextureSlot = gl.TEXTURE2
	BatchTextureSlotProjectionsIdxBuffer BatchTextureSlot = gl.TEXTURE3
	BatchTextureSlotViewsIdxBuffer       BatchTextureSlot = gl.TEXTURE4
)

type TextureSlot uint32

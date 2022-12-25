package render

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// TODO: collect views and projections
// assigned to cameras and canvases of
// the batched sprites. Pass a view
// and a projection index for each sprite.

// TODO: add a shader program field to batch.
// The projections and views should be
// uniform arrays of predefined sizes.
// The user can assign the initial size
// to the shader program, and if he needs
// to go beyond it, the shader program must
// be recompiled. The absolute max uniform
// array size cannot go past 256 because the
// projection index number for a vertice
// is only 1 byte long. Everytime the user
// adds a new canvas to the layout, the batch
// shader program uniforms are reassigned.

// TODO: handle elements.

// TODO: use sampler buffers for models,
// should draw flags, projection and view indices.

type Batch struct {
	glHandler                            uint32 // glHandler is an OpenGL name for the underlying batch VAO.
	glModelsTextureBufferHandler         uint32
	glShouldDrawTextureBufferHandler     uint32
	glProjectionsIdxTextureBufferHandler uint32
	glViewsIdxTextureBufferHandler       uint32

	sprites       []*Sprite
	layout        *Layout
	texture       *Texture
	shaderProgram *ShaderProgram

	viewsBuffer       []mgl32.Mat4
	projectionsBuffer []mgl32.Mat4

	// TODO: everytime we change a
	// parameter of the sprite we
	// should also change it in the
	// corresponding GPU list by the
	// batch index of the sprite.
	projectionsIdx *gpuList[byte]
	models         *gpuList[float32]
	viewsIdx       *gpuList[byte]
	vertices       *gpuList[float32]
	texCoords      *gpuList[float32]
	colorMasks     *gpuList[float32]
	shouldDraw     *gpuList[byte]
}

func (batch *Batch) rebindModelsTextureBuffer() {
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_BUFFER, batch.glModelsTextureBufferHandler)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, batch.models.glHandler)

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (batch *Batch) rebindShouldDrawTextureBuffer() {
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_BUFFER, batch.glShouldDrawTextureBufferHandler)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, batch.shouldDraw.glHandler)

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (batch *Batch) rebindProjectionsIdxTextureBuffer() {
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_BUFFER, batch.glProjectionsIdxTextureBufferHandler)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R8, batch.projectionsIdx.glHandler)

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (batch *Batch) rebindViewsIdxTextureBuffer() {
	gl.ActiveTexture(gl.TEXTURE4)
	gl.BindTexture(gl.TEXTURE_BUFFER, batch.glViewsIdxTextureBufferHandler)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R8, batch.viewsIdx.glHandler)

	gl.ActiveTexture(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (batch *Batch) setCanvasProjection(idx int, projection mgl32.Mat4) {
	batch.shaderProgram.Use()
	defer gl.UseProgram(0)

	batch.projectionsBuffer[idx] = projection

	header := *(*reflect.SliceHeader)(unsafe.Pointer(&batch.projectionsBuffer))
	header.Len *= 16
	header.Cap *= 16
	data := *(*[]float32)(unsafe.Pointer(&header))

	batch.shaderProgram.SetFloat32Array("projections", data)
}

func (batch *Batch) setCanvasView(idx int, view mgl32.Mat4) {
	batch.shaderProgram.Use()
	defer gl.UseProgram(0)

	batch.viewsBuffer[idx] = view
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&batch.viewsBuffer))
	header.Len *= 16
	header.Cap *= 16
	data := *(*[]float32)(unsafe.Pointer(&header))

	batch.shaderProgram.SetFloat32Array("views", data)
}

// TODO: bind all the texture buffers
// as textures in slots in Draw().

func (batch *Batch) Draw() {
	batch.shaderProgram.Use()
	defer gl.UseProgram(0)

	gl.BindVertexArray(batch.glHandler)
	defer gl.BindVertexArray(0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, spriteIndexBufferHandler)
	defer gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	batch.texture.Use()
	defer func() {
		gl.ActiveTexture(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}()

	// Send all the canvas views to the GPU.
	for i := 0; i < len(batch.layout.canvases); i++ {
		canvas := batch.layout.canvases[i]
		batch.viewsBuffer[i] = canvas.camera.View()
	}

	header := *(*reflect.SliceHeader)(unsafe.Pointer(&batch.viewsBuffer))
	header.Len *= 16
	header.Cap *= 16
	data := *(*[]float32)(unsafe.Pointer(&header))

	batch.shaderProgram.SetFloat32Array("views", data)

	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_INT,
		gl.PtrOffset(0), int32(len(batch.sprites)))
}

func (batch *Batch) AttachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	if sprite.texture != batch.texture {
		return fmt.Errorf(
			"the sprite should have the same texture as the batch")
	}

	ind := len(batch.sprites)
	sprite.batchIndex = ind
	sprite.batch = batch
	batch.sprites = append(batch.sprites, sprite)

	identMatrix := mgl32.Ident4()

	batch.vertices.addDataFromBuffer(
		sprite.glVertexBufferHandler, 8)
	batch.texCoords.addDataFromBuffer(
		sprite.glTextureCoordinatesBufferHandler, 8)
	batch.colorMasks.addDataFromBuffer(
		sprite.glColorMaskBufferHandler, 16)

	prevCapacity := batch.projectionsIdx.getCapacity()
	batch.projectionsIdx.addElement(byte(sprite.canvas.index))

	if batch.projectionsIdx.getCapacity() > prevCapacity {
		batch.rebindProjectionsIdxTextureBuffer()
	}

	prevCapacity = batch.viewsIdx.getCapacity()
	batch.viewsIdx.addElement(byte(sprite.canvas.index))

	if batch.viewsIdx.getCapacity() > prevCapacity {
		batch.rebindViewsIdxTextureBuffer()
	}

	prevCapacity = batch.models.getCapacity()
	batch.models.addElements(identMatrix[:])

	if batch.models.getCapacity() > prevCapacity {
		batch.rebindModelsTextureBuffer()
	}

	prevCapacity = batch.shouldDraw.getCapacity()
	batch.shouldDraw.addElement(0)

	if batch.shouldDraw.getCapacity() > prevCapacity {
		batch.rebindShouldDrawTextureBuffer()
	}

	sprite.deleteVertexBuffer()
	sprite.deleteTextureCoordinatesBuffer()
	sprite.deleteColorMaskBuffer()
	sprite.deleteVertexArray()

	return nil
}

func (batch *Batch) DetachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	length := len(batch.sprites)
	ind := sprite.batchIndex

	if ind < length && batch.sprites[ind] != sprite ||
		ind >= length || ind < 0 {
		return fmt.Errorf(
			"the sprite doesn't exist on the batch")
	}

	// Reindex all the sprites
	// remaining on the batch.
	for i := ind + 1; i < length; i++ {
		batch.sprites[i].batchIndex--
	}

	if ind == 0 {
		batch.sprites = batch.sprites[1:]
	} else if ind < length-1 {
		batch.sprites = append(batch.sprites[:ind],
			batch.sprites[ind+1:]...)
	} else {
		batch.sprites = batch.sprites[:length-1]
	}

	sprite.batch = nil
	batchIndex := sprite.batchIndex
	sprite.batchIndex = -1

	sprite.createVertexBuffer()
	sprite.createTextureCoordinatesBuffer()
	sprite.createColorMaskBuffer()

	batch.vertices.copyDataToBuffer(sprite.
		glVertexBufferHandler, batchIndex, 8)
	batch.texCoords.copyDataToBuffer(sprite.
		glTextureCoordinatesBufferHandler, batchIndex, 8)
	batch.colorMasks.copyDataToBuffer(sprite.
		glColorMaskBufferHandler, batchIndex, 16)

	sprite.createVertexArray()
	sprite.assembleVertexArray()

	err := batch.vertices.removeElements(batchIndex, 8)

	if err != nil {
		return err
	}

	err = batch.texCoords.removeElements(batchIndex, 8)

	if err != nil {
		return err
	}

	err = batch.colorMasks.removeElements(batchIndex, 16)

	if err != nil {
		return err
	}

	err = batch.projectionsIdx.removeElement(batchIndex)

	if err != nil {
		return err
	}

	err = batch.viewsIdx.removeElement(batchIndex)

	if err != nil {
		return err
	}

	err = batch.models.removeElements(batchIndex, 16)

	if err != nil {
		return err
	}

	err = batch.shouldDraw.removeElement(batchIndex)

	if err != nil {
		return err
	}

	return nil
}

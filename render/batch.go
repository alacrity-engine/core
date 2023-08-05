package render

import (
	_ "embed"
	"fmt"
	"sort"

	"github.com/alacrity-engine/core/geometry"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	//go:embed std-batch-frag.glsl
	batchFragmentShaderSource string
	//go:embed std-batch-vert.glsl
	batchVertexShaderSource string
)

// TODO: a batch should have its own canvas
// and accept only the sprites that belong to it
// because canvases are drawn sequentially and
// they may have other sprites besides batched ones.

// TODO: create a cache for frequently used runtime
// objects (for example, compiled standard shaders and programs)
// so there is no need to create the same object many times.
// Or maybe just create a global variable for each one of them.

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

// TODO: compute maxNumCanvases from the
// current number of canvases on the layout.

type Batch struct {
	glHandler                   uint32 // glHandler is an OpenGL name for the underlying batch VAO.
	modelsTextureBuffer         *TextureBuffer
	shouldDrawTextureBuffer     *TextureBuffer
	projectionsIdxTextureBuffer *TextureBuffer
	viewsIdxTextureBuffer       *TextureBuffer

	sprites       []*Sprite
	canvas        *Canvas
	texture       *Texture
	shaderProgram *ShaderProgram

	models     *gpuList[float32]
	vertices   *gpuList[float32]
	texCoords  *gpuList[float32]
	colorMasks *gpuList[float32]
	shouldDraw *gpuList[byte]
}

func (batch *Batch) buildVAO() {
	gl.BindVertexArray(batch.glHandler)
	defer gl.BindVertexArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vertices.glHandler)
	vertAttrib := uint32(gl.GetAttribLocation(batch.shaderProgram.glHandler, gl.Str("aPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, batch.texCoords.glHandler)
	texCoordAttrib := uint32(gl.GetAttribLocation(batch.shaderProgram.glHandler, gl.Str("aTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, batch.colorMasks.glHandler)
	colorAttrib := uint32(gl.GetAttribLocation(batch.shaderProgram.glHandler, gl.Str("aColor\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
}

func (batch *Batch) Draw() {
	//gl.Disable(gl.DEPTH_TEST)
	//defer gl.Enable(gl.DEPTH_TEST)

	batch.shaderProgram.Use()
	defer gl.UseProgram(0)

	gl.BindVertexArray(batch.glHandler)
	defer gl.BindVertexArray(0)

	batch.texture.Use()
	batch.modelsTextureBuffer.Bind()
	batch.shouldDrawTextureBuffer.Bind()
	batch.projectionsIdxTextureBuffer.Bind()
	batch.viewsIdxTextureBuffer.Bind()

	defer func() {
		gl.ActiveTexture(0)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.BindTexture(gl.TEXTURE_BUFFER, 0)
	}()

	batch.shaderProgram.SetMatrix4("view", batch.canvas.camera.View())
	batch.shaderProgram.SetMatrix4("projection", batch.canvas.projection)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(batch.sprites)*6))

	batch.shouldDraw.clear()
}

func (batch *Batch) AttachSprite(sprite *Sprite) error {
	if sprite == nil {
		return fmt.Errorf("the sprite is nil")
	}

	if sprite.texture != batch.texture {
		return fmt.Errorf(
			"the sprite should have the same texture as the batch")
	}

	ind := sort.Search(len(batch.sprites), func(i int) bool {
		return batch.sprites[i].drawZ >= sprite.drawZ
	})

	sprite.batchIndex = ind
	sprite.batch = batch

	if ind <= 0 {
		batch.sprites = append([]*Sprite{sprite}, batch.sprites...)
	} else if ind >= len(batch.sprites) {
		batch.sprites = append(batch.sprites, sprite)
	} else {
		batch.sprites = append(batch.sprites[:ind+1], batch.sprites[ind:]...)
		batch.sprites[ind] = sprite
	}

	// Reindex all the sprites
	// remaining on the batch.
	for i := ind + 1; i < len(batch.sprites); i++ {
		batch.sprites[i].batchIndex++
	}

	vertices := make([]float32, 18)
	geometry.ComputeSpriteVerticesNoElementsFill(
		vertices, width, height, sprite.targetArea)
	err := batch.vertices.insertElements(
		ind*len(vertices), len(vertices), vertices)

	if err != nil {
		return err
	}

	texCoords := make([]float32, 12)
	geometry.ComputeSpriteTextureCoordinatesNoElementsFill(
		texCoords, sprite.texture.imageWidth,
		sprite.texture.imageHeight, sprite.targetArea)
	err = batch.texCoords.insertElements(
		ind*len(texCoords), len(texCoords), texCoords)

	if err != nil {
		return err
	}

	colorMask := make([]float32, 24)
	geometry.ColorMaskDataNoElementsFill(
		colorMask, sprite.colorMask.Data())
	err = batch.colorMasks.insertElements(
		ind*len(colorMask), len(colorMask), colorMask)

	if err != nil {
		return err
	}

	// Check the proportion.
	//a := batch.vertices.capacity / batch.vertices.stride
	//b := batch.texCoords.capacity / batch.texCoords.stride
	//c := batch.colorMasks.capacity / batch.colorMasks.stride
	//
	//if a != b && b != c && a != c {
	//	_ = a
	//}

	// Rebind all the sprite data to the VAO.
	batch.buildVAO()

	prevCapacity := batch.models.getCapacity()
	identMatrix := mgl32.Ident4()
	err = batch.models.insertElements(
		ind*len(identMatrix), len(identMatrix), identMatrix[:])

	if err != nil {
		return err
	}

	if batch.models.getCapacity() > prevCapacity {
		batch.modelsTextureBuffer.
			rebind(batch.models.glHandler)
	}

	prevCapacity = batch.shouldDraw.getCapacity()
	batch.shouldDraw.insertElement(ind, 0)

	if batch.shouldDraw.getCapacity() > prevCapacity {
		batch.shouldDrawTextureBuffer.
			rebind(batch.shouldDraw.glHandler)
	}

	sprite.deleteVertexBuffer()
	sprite.deleteTextureCoordinatesBuffer()
	sprite.deleteColorMaskBuffer()
	sprite.deleteVertexArray()
	err = sprite.canvas.removeBatchedSprite(sprite)

	if err != nil {
		return err
	}

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
	for i := ind; i < length; i++ {
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

	err := sprite.SetTargetArea(sprite.targetArea)

	if err != nil {
		return err
	}

	err = sprite.SetColorMask(sprite.colorMask)

	if err != nil {
		return err
	}

	sprite.createVertexArray()
	sprite.assembleVertexArray()

	err = batch.vertices.removeElements(batchIndex, 18)

	if err != nil {
		return err
	}

	err = batch.texCoords.removeElements(batchIndex, 12)

	if err != nil {
		return err
	}

	err = batch.colorMasks.removeElements(batchIndex, 24)

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

	err = sprite.canvas.addSpriteFromBatch(sprite)

	if err != nil {
		return err
	}

	return nil
}

func NewBatch(texture *Texture, canvas *Canvas, options ...BatchOption) (*Batch, error) {
	params := batchParameters{initialObjectCapacity: 0}
	var batch Batch

	for i := 0; i < len(options); i++ {
		option := options[i]
		err := option(&batch, &params)

		if err != nil {
			return nil, err
		}
	}

	// Instantiate GPU lists.
	batch.models = newGPUList(DrawModeStream,
		make([]float32, params.initialObjectCapacity), 1)
	batch.vertices = newGPUList(DrawModeDynamic,
		make([]float32, params.initialObjectCapacity), 3*4)
	batch.texCoords = newGPUList(DrawModeDynamic,
		make([]float32, params.initialObjectCapacity), 2*4)
	batch.colorMasks = newGPUList(DrawModeDynamic,
		make([]float32, params.initialObjectCapacity), 4*4)
	batch.shouldDraw = newGPUList(DrawModeStream,
		make([]byte, params.initialObjectCapacity), 1)

	// Instantiate texture buffers.
	batch.modelsTextureBuffer = NewTextureBuffer(batch.models.glHandler,
		TextureSlot(BatchTextureSlotModelsBuffer), FormatFloat32)
	batch.shouldDrawTextureBuffer = NewTextureBuffer(batch.models.glHandler,
		TextureSlot(BatchTextureSlotShouldDrawBuffer), FormatByte)
	batch.projectionsIdxTextureBuffer = NewTextureBuffer(batch.models.glHandler,
		TextureSlot(BatchTextureSlotProjectionsIdxBuffer), FormatByte)
	batch.viewsIdxTextureBuffer = NewTextureBuffer(batch.models.glHandler,
		TextureSlot(BatchTextureSlotViewsIdxBuffer), FormatByte)

	batch.sprites = make([]*Sprite, 0, params.initialObjectCapacity)
	batch.texture = texture
	batch.canvas = canvas

	// Compile batch shaders.
	vertexShader, _, err := NewStandardBatchShader(
		ShaderTypeVertex)

	if err != nil {
		return nil, err
	}

	fragmentShader, _, err := NewStandardBatchShader(
		ShaderTypeFragment)

	if err != nil {
		return nil, err
	}

	batch.shaderProgram, err = NewShaderProgramFromShaders(
		vertexShader, fragmentShader)

	if err != nil {
		return nil, err
	}

	// Build a VAO.
	var handler uint32
	gl.GenVertexArrays(1, &handler)
	batch.glHandler = handler
	batch.buildVAO()

	return &batch, nil
}

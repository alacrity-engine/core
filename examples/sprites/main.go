package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"runtime"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/system"
	"github.com/alacrity-engine/core/system/collections"
	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/image/colornames"
)

const (
	width  = 1920
	height = 1080
)

func init() {
	runtime.LockOSThread()
}

// TODO: add an opportunity to assign a size and a transform to a canvas.
// I think I'll need to make all the canvas sprites' transforms children
// of the canvas transform in order to move the canvas with its contents.
// This will allow me make different versions of the same game for 4:3 and
// 16:9 aspects.

// TODO: add batch rendering for static and animated sprites.

func main() {
	// Initialize the engine.
	err := system.InitializeWindow("Demo", width, height, false, false)
	handleError(err)
	err = render.Initialize(width, height, -30, 30)
	handleError(err)

	// Initialize the shader program.
	shaderProgram, err := render.NewStandardSpriteShaderProgram()
	handleError(err)

	// Load the Cirno picture.
	file, err := os.Open("cirno.png")
	handleError(err)
	img, _, err := image.Decode(file)
	handleError(err)
	imgRGBA := img.(*image.RGBA)
	reversePix(imgRGBA.Pix)
	mirror(imgRGBA)
	err = file.Close()
	handleError(err)

	// Create a texture and a sprite for Cirno.
	cirnoTexture := render.NewTextureFromImage(imgRGBA, render.TextureFilteringLinear)
	cirnoSprite, err := render.NewSpriteFromTextureAndProgram(
		render.DrawModeStatic, render.DrawModeStatic,
		render.DrawModeStatic, cirnoTexture, shaderProgram,
		geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)
	cirnoSprite.SetColorMask(render.RGBARepeat4(render.ToRGBA(colornames.Red)))

	// Load the Sakuya picture.
	file, err = os.Open("sakuya.png")
	handleError(err)
	img, _, err = image.Decode(file)
	handleError(err)
	imgRGBA = img.(*image.RGBA)
	reversePix(imgRGBA.Pix)
	mirror(imgRGBA)
	err = file.Close()
	handleError(err)

	// Create a texture and a sprite for Sakuya.
	sakuyaTexture := render.NewTextureFromImage(imgRGBA, render.TextureFilteringLinear)
	sakuyaSprite, err := render.NewSpriteFromTextureAndProgram(
		render.DrawModeStatic, render.DrawModeStatic,
		render.DrawModeStatic, sakuyaTexture, shaderProgram,
		geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)
	sakuyaSprite.SetColorMask(render.ColorMask{
		render.ToRGBA(colornames.Red),
		render.ToRGBA(colornames.Green),
		render.ToRGBA(colornames.Blue),
		render.ToRGBA(colornames.Black),
	})

	// Add canvases.
	layout := render.NewLayout()

	zBufferDataNodePool, err := mempool.NewPool[*avltree.AVLNode[int64, *render.Sprite]](
		func() *avltree.AVLNode[int64, *render.Sprite] {
			return new(avltree.AVLNode[int64, *render.Sprite])
		})
	handleError(err)
	zBufferNodePool, err := mempool.NewPool[*collections.UnrestrictedAVLNode[render.Geometric, render.ZBufferData]](
		func() *collections.UnrestrictedAVLNode[render.Geometric, render.ZBufferData] {
			return new(collections.UnrestrictedAVLNode[render.Geometric, render.ZBufferData])
		})
	handleError(err)
	zBufferDataPool, err := mempool.NewPool[*collections.AVLSortedDictionary[int64, *render.Sprite]](
		func() *collections.AVLSortedDictionary[int64, *render.Sprite] {
			tree, _ := collections.NewAVLSortedDictionary[int64, *render.Sprite]()
			return tree
		},
	)
	handleError(err)
	zBufferPool, err := mempool.NewPool[*collections.AVLUnrestrictedSortedDictionary[render.Geometric, render.ZBufferData]](
		func() *collections.AVLUnrestrictedSortedDictionary[render.Geometric, render.ZBufferData] {
			tree, _ := collections.NewAVLUnrestrictedSortedDictionary[render.Geometric, render.ZBufferData]()
			return tree
		},
	)
	handleError(err)

	zBufferDataProducer := collections.NewAVLSortedDictionaryProducer[int64, *render.Sprite](
		zBufferDataPool, zBufferDataNodePool)
	zBufferProducer := collections.NewAVLUnrestrictedSortedDictionaryProducer[render.Geometric, render.ZBufferData](
		zBufferPool, zBufferNodePool)

	cirnoCanvas, err := render.NewCanvas("cirno", 4, render.Ortho2DStandard(),
		zBufferProducer, zBufferDataProducer)
	handleError(err)
	err = layout.AddCanvas(cirnoCanvas)
	handleError(err)
	err = cirnoCanvas.AddSprite(cirnoSprite)
	handleError(err)

	sakuyaCanvas, err := render.NewCanvas("sakuya", 2, render.Ortho2DStandard(),
		zBufferProducer, zBufferDataProducer)
	handleError(err)
	err = layout.AddCanvas(sakuyaCanvas)
	handleError(err)
	err = sakuyaCanvas.AddSprite(sakuyaSprite)
	handleError(err)

	sakuyaTransform := geometry.NewTransform(nil)
	cirnoTransform := geometry.NewTransform(nil)
	//transform.ApplyScale(geometry.V(0.5, 0.5))

	system.InitMetrics()

	fmt.Println("GPU vendor:", system.Vendor())
	fmt.Println("GPU model:", system.Renderer())

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		//moveCamera(sakuyaCanvas.Camera(), 1000, system.DeltaTime())

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)
		sakuyaTransform.Rotate(math.Pi / 4.0 * geometry.RadToDeg * system.DeltaTime())
		sakuyaCanvas.Camera().Move(geometry.V(200, 200).Scaled(system.DeltaTime()))
		sakuyaTransform.Move(geometry.V(200, 200).Scaled(system.DeltaTime()))
		cirnoTransform.Move(geometry.V(400, 400).Scaled(system.DeltaTime()))
		cirnoSprite.Draw(cirnoTransform)
		sakuyaSprite.Draw(sakuyaTransform)

		err = layout.Draw()
		handleError(err)

		system.TickLoop()
		system.UpdateFrameRate()
	}
}

func moveCamera(camera *render.Camera, speed, deltaTime float64) {
	movement := geometry.ZV

	if system.ButtonPressed(system.KeyUp) {
		movement.Y += 1.0
	}

	if system.ButtonPressed(system.KeyDown) {
		movement.Y -= 1.0
	}

	if system.ButtonPressed(system.KeyLeft) {
		movement.X -= 1.0
	}

	if system.ButtonPressed(system.KeyRight) {
		movement.X += 1.0
	}

	if movement != geometry.ZV {
		movement = geometry.ClampMagnitude(movement, 1.0)
		movement = movement.Scaled(speed * deltaTime)

		camera.Move(movement)
	}
}

func reversePix(arr []byte) {
	start := 0
	end := len(arr) - 4

	for start < end {
		for i := 0; i < 4; i++ {
			temp := arr[start+i]
			arr[start+i] = arr[end+i]
			arr[end+i] = temp
		}

		start += 4
		end -= 4
	}
}

func mirror(img *image.RGBA) {
	for i := 0; i < img.Rect.Dy(); i++ {
		reversePix(img.Pix[4*i*img.Rect.Dx() : 4*(i+1)*img.Rect.Dx()])
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

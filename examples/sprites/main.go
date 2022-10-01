package main

import (
	"image"
	_ "image/png"
	"math"
	"os"
	"runtime"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/system"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

const (
	width  = 1920
	height = 1080
)

func init() {
	runtime.LockOSThread()
}

// TODO: add a gradient (multi-color) color mask.

// TODO: add batch rendering.

// TODO: add a 2D camera type and assign a camera
// to each canvas.

func main() {
	// Initialize the engine.
	err := system.InitializeWindow("Demo", width, height, false, false)
	handleError(err)
	err = render.Initialize(width, height)
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
	cirnoSprite, err := render.NewSpriteFromTextureAndProgram(render.DrawModeStatic,
		cirnoTexture, shaderProgram, geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)

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
	sakuyaSprite, err := render.NewSpriteFromTextureAndProgram(render.DrawModeStatic,
		sakuyaTexture, shaderProgram, geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)

	// Add canvases.
	layout := render.NewLayout()

	cirnoCanvas := render.NewCanvas(0)
	layout.AddCanvas(cirnoCanvas)
	cirnoCanvas.AddSprite(cirnoSprite)

	sakuyaCanvas := render.NewCanvas(2)
	layout.AddCanvas(sakuyaCanvas)
	sakuyaCanvas.AddSprite(sakuyaSprite)

	aspect := float32(height) / float32(width)
	projection := mgl32.Ortho(-1, 1, -1*aspect, 1*aspect, -1, 1)
	sakuyaTransform := geometry.NewTransform(nil)
	cirnoTransform := geometry.NewTransform(nil)
	//transform.ApplyScale(geometry.V(0.5, 0.5))
	//transform.Move(geometry.V(0.5, 0.5))

	system.InitMetrics()

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)
		sakuyaTransform.Rotate(math.Pi / 4.0 * geometry.RadToDeg * system.DeltaTime())
		sakuyaTransform.Move(geometry.V(200, 200).Scaled(system.DeltaTime()))
		cirnoTransform.Move(geometry.V(400, 400).Scaled(system.DeltaTime()))
		cirnoSprite.Draw(cirnoTransform.Data(), mgl32.Ident4(), projection)
		sakuyaSprite.Draw(sakuyaTransform.Data(), mgl32.Ident4(), projection)

		system.TickLoop()
		system.UpdateFrameRate()
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

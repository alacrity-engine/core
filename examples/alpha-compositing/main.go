package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
	"runtime"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/system"
	"golang.org/x/image/colornames"
)

const (
	width  = 1920
	height = 1080
)

func init() {
	runtime.LockOSThread()
}

func main() {
	// Initialize the engine.
	err := system.InitializeWindow("Demo", width, height, false, false)
	handleError(err)
	err = render.Initialize(width, height, -30, 30)
	handleError(err)

	// Load the ball picture.
	file, err := os.Open("ball.png")
	handleError(err)
	img, _, err := image.Decode(file)
	handleError(err)
	imgRGBA := image.NewRGBA(image.Rect(
		0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(imgRGBA, imgRGBA.Bounds(),
		img, img.Bounds().Min, draw.Src)
	reversePix(imgRGBA.Pix)
	mirror(imgRGBA)
	err = file.Close()
	handleError(err)

	// Initialize the shader program.
	shaderProgram, err := render.NewStandardSpriteShaderProgram()
	handleError(err)

	// Create a texture and a sprite for the ball.
	ballTexture := render.NewTextureFromImage(imgRGBA, render.TextureFilteringLinear)
	ballSprite, err := render.NewSpriteFromTextureAndProgram(
		render.DrawModeStatic, render.DrawModeStatic,
		render.DrawModeStatic, ballTexture, shaderProgram,
		geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)

	// Add canvases.
	layout := render.NewLayout()

	ballCanvas := render.NewCanvas(0, render.Ortho2DStandard())
	err = layout.AddCanvas(ballCanvas)
	handleError(err)
	err = ballCanvas.AddSprite(ballSprite)
	handleError(err)

	ballTransform := geometry.NewTransform(nil)

	system.InitMetrics()

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)

		ballSprite.SetZ(2)
		ballSprite.Draw(ballTransform)
		ballSprite.SetZ(0)
		ballTransform.MoveTo(geometry.V(float64(imgRGBA.Bounds().Dx()/2), float64(imgRGBA.Bounds().Dy()/2)))
		ballSprite.Draw(ballTransform)
		ballTransform.MoveTo(geometry.V(-float64(imgRGBA.Bounds().Dx()/2), -float64(imgRGBA.Bounds().Dy()/2)))
		ballSprite.SetZ(2)

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

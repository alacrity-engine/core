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

// TODO: resolve the gl.DrawElemnts SIGSEGV.

// TODO: add batch rendering for static and animated sprites.

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
	cirnoSprite, err := render.NewSpriteFromTextureAndProgram(render.DrawModeStatic, render.DrawModeStatic,
		cirnoTexture, shaderProgram, geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
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
	sakuyaSprite, err := render.NewSpriteFromTextureAndProgram(render.DrawModeStatic, render.DrawModeStatic,
		sakuyaTexture, shaderProgram, geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)
	sakuyaSprite.SetColorMask([4]render.RGBA{
		render.ToRGBA(colornames.Red),
		render.ToRGBA(colornames.Green),
		render.ToRGBA(colornames.Blue),
		render.ToRGBA(colornames.Black),
	})

	// Add canvases.
	layout := render.NewLayout()

	cirnoCanvas := render.NewCanvas(0)
	layout.AddCanvas(cirnoCanvas)
	cirnoCanvas.AddSprite(cirnoSprite)

	sakuyaCanvas := render.NewCanvas(2)
	layout.AddCanvas(sakuyaCanvas)
	sakuyaCanvas.AddSprite(sakuyaSprite)

	sakuyaTransform := geometry.NewTransform(nil)
	cirnoTransform := geometry.NewTransform(nil)
	//transform.ApplyScale(geometry.V(0.5, 0.5))

	system.InitMetrics()

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
		cirnoSprite.DrawTransform(cirnoTransform)
		sakuyaSprite.DrawTransform(sakuyaTransform)

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

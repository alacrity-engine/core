package main

import (
	"image"
	_ "image/png"
	"os"
	"runtime"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/system"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/colornames"
)

const (
	width  = 800
	height = 600
)

func init() {
	runtime.LockOSThread()
}

// TODO: the sprite always reaches the upper
// edge of the window at the same time regardless
// of the window resoltion thus moving with
// different speed that depends on the screen size.
// This should be fixed.

// TODO: add a canvas (a group of sprites with
// the same Z drawing coordinate).

func main() {
	file, err := os.Open("cirno.png")
	handleError(err)
	img, _, err := image.Decode(file)
	handleError(err)
	imgRGBA := img.(*image.RGBA)
	reversePix(imgRGBA.Pix)
	mirror(imgRGBA)

	err = system.InitializeWindow("Demo", width, height, false, false)
	handleError(err)
	err = render.Initialize(width, height)
	handleError(err)

	shaderProgram, err := render.NewStandardSpriteShaderProgram()
	handleError(err)

	texture := render.NewTextureFromImage(imgRGBA, render.TextureFilteringLinear)
	sprite, err := render.NewSpriteFromTextureAndProgram(render.DrawModeStatic,
		texture, shaderProgram, geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)

	aspect := float32(height) / float32(width)
	projection := mgl32.Ortho(-1, 1, -1*aspect, 1*aspect, -1, 1)
	//model := mgl32.Ident4()
	//model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0))
	transform := geometry.NewTransform(nil)
	//transform.ApplyScale(geometry.V(0.5, 0.5))

	system.InitMetrics()

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)
		//model = model.Mul4(mgl32.Scale3D(1.1*float32(system.DeltaTime()),
		//	1.1*float32(system.DeltaTime()), 0))
		//model = model.Mul4(mgl32.HomogRotate3DZ(float32(math.Pi/4.0) * float32(
		//	system.DeltaTime())))
		//model = mgl32.Translate3D(0.1*float32(system.DeltaTime()),
		//	0.1*float32(system.DeltaTime()), 0).Mul4(model)
		//transform.Rotate(math.Pi / 4.0 * geometry.RadToDeg * system.DeltaTime())
		transform.Move(geometry.V(0.1, 0.1).Scaled(system.DeltaTime()))
		sprite.Draw(transform.Data(), mgl32.Ident4(), projection)

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

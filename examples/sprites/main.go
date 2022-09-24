package main

import (
	"image"
	_ "image/png"
	"io/ioutil"
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

func main() {
	file, err := os.Open("sakuya.png")
	handleError(err)
	img, _, err := image.Decode(file)
	handleError(err)
	imgRGBA := img.(*image.RGBA)
	reversePix(imgRGBA.Pix)
	mirror(imgRGBA)

	vertexShaderSourceData, err := ioutil.ReadFile("vert.glsl")
	handleError(err)
	vertexShaderSource := string(vertexShaderSourceData) + "\x00"
	fragmentShaderSourceData, err := ioutil.ReadFile("frag.glsl")
	handleError(err)
	fragmentShaderSource := string(fragmentShaderSourceData) + "\x00"

	err = system.InitializeWindow("Demo", width, height, false, false)
	handleError(err)
	err = render.Initialize(width, height)
	handleError(err)

	vertexShader, err := render.NewShaderFromSource(
		vertexShaderSource, render.ShaderTypeVertex)
	handleError(err)
	fragmentShader, err := render.NewShaderFromSource(
		fragmentShaderSource, render.ShaderTypeFragment)
	handleError(err)
	shaderProgram, err := render.NewShaderProgramFromShaders(
		vertexShader, fragmentShader)
	handleError(err)

	texture := render.NewTextureFromImage(imgRGBA, render.TextureFilteringLinear)
	sprite, err := render.NewSpriteFromTextureAndProgram(render.DrawModeStatic,
		texture, shaderProgram, geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)

	aspect := float32(height) / float32(width)
	projection := mgl32.Ortho(-1, 1, -1*aspect, 1*aspect, -1, 1)

	for !system.ShouldClose() {
		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)
		sprite.Draw(mgl32.Ident4(), mgl32.Ident4(), projection)
		system.TickLoop()
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

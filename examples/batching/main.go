package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
	"runtime"

	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/system"
	"github.com/alacrity-engine/core/system/collections"
	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/image/colornames"
)

type Ball struct {
	Sprite    *render.Sprite
	Transform *geometry.Transform
}

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

	ballCanvas, err := render.NewCanvas("balls", 0, render.Ortho2DStandard(), zBufferProducer, zBufferDataProducer)
	handleError(err)
	err = layout.AddCanvas(ballCanvas)
	handleError(err)

	// Create batch.
	batch, err := render.NewBatch("balls", ballTexture)
	handleError(err)
	err = ballCanvas.AddBatch(batch, -15, 15)
	handleError(err)
	balls := make([]*Ball, 0, width*height/
		(imgRGBA.Bounds().Dx()*imgRGBA.Bounds().Dy()))

	// Instantiate all the objects and
	// attach them to the batch.
	/*zCounter := 0

	for i := 0; i < 32; i++ {
		for j := 0; j < 32; j++ {
			x := float64(i) / 32.0 * width
			y := float64(j) / 32.0 * height

			x -= width / 2
			y -= height / 2

			ballSprite, err := render.NewSpriteFromTextureAndProgram(
				render.DrawModeStatic, render.DrawModeStatic,
				render.DrawModeStatic, ballTexture, shaderProgram,
				geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
			handleError(err)
			ballSprite.SetZ(float32(zCounter))
			err = ballCanvas.AttachSpriteToBatch(batch, ballSprite)
			handleError(err)
			ballTransform := geometry.NewTransform(nil)
			ballTransform.MoveTo(geometry.V(x, y))

			ball := &Ball{
				Sprite:    ballSprite,
				Transform: ballTransform,
			}
			balls = append(balls, ball)

			zCounter++
		}
	}*/

	ballSprite1, err := render.NewSpriteFromTextureAndProgram(
		render.DrawModeStatic, render.DrawModeStatic,
		render.DrawModeStatic, ballTexture, shaderProgram,
		geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)
	ballSprite2, err := render.NewSpriteFromTextureAndProgram(
		render.DrawModeStatic, render.DrawModeStatic,
		render.DrawModeStatic, ballTexture, shaderProgram,
		geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)
	ballSprite3, err := render.NewSpriteFromTextureAndProgram(
		render.DrawModeStatic, render.DrawModeStatic,
		render.DrawModeStatic, ballTexture, shaderProgram,
		geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
	handleError(err)

	ballTransform1 := geometry.NewTransform(nil)
	ballTransform2 := geometry.NewTransform(nil)
	ballTransform3 := geometry.NewTransform(nil)

	// Upper right.
	ballSprite1.SetZ(-1)
	ballTransform1.MoveTo(geometry.V(float64(imgRGBA.Bounds().Dx()/4), float64(imgRGBA.Bounds().Dy()/4)))

	// Lower left.
	ballSprite2.SetZ(1)
	ballTransform2.MoveTo(geometry.V(-float64(imgRGBA.Bounds().Dx()/4), -float64(imgRGBA.Bounds().Dy()/4)))

	err = ballCanvas.AttachSpriteToBatch(batch, ballSprite1)
	handleError(err)
	err = ballCanvas.AttachSpriteToBatch(batch, ballSprite2)
	handleError(err)
	err = ballCanvas.AttachSpriteToBatch(batch, ballSprite3)
	handleError(err)

	balls = append(balls, &Ball{
		Sprite:    ballSprite1,
		Transform: ballTransform1,
	}, &Ball{
		Sprite:    ballSprite2,
		Transform: ballTransform2,
	}, &Ball{
		Sprite:    ballSprite3,
		Transform: ballTransform3,
	})

	var upperRemoved, middleRemoved, lowerRemoved bool

	// Load the Cirno picture.
	/*file, err = os.Open("cirno.png")
	handleError(err)
	img, _, err = image.Decode(file)
	handleError(err)
	imgRGBA = img.(*image.RGBA)
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
	cirnoTransform := geometry.NewTransform(nil)

	cirnoSprite.SetZ(16)
	err = ballCanvas.AddSprite(cirnoSprite)
	handleError(err)*/

	system.InitMetrics()

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		if system.ButtonPressed(system.KeyA) && !upperRemoved {
			err = batch.DetachSprite(ballSprite1)
			handleError(err)

			balls = []*Ball{
				{
					Sprite:    ballSprite2,
					Transform: ballTransform2,
				},
				{
					Sprite:    ballSprite3,
					Transform: ballTransform3,
				},
			}

			upperRemoved = true
		}

		if system.ButtonPressed(system.KeyD) && !lowerRemoved {
			err = batch.DetachSprite(ballSprite2)
			handleError(err)

			balls = []*Ball{
				{
					Sprite:    ballSprite1,
					Transform: ballTransform1,
				},
				{
					Sprite:    ballSprite3,
					Transform: ballTransform3,
				},
			}

			lowerRemoved = true
		}

		if system.ButtonPressed(system.KeyS) && !middleRemoved {
			err = batch.DetachSprite(ballSprite3)
			handleError(err)

			balls = []*Ball{
				{
					Sprite:    ballSprite1,
					Transform: ballTransform1,
				},
				{
					Sprite:    ballSprite2,
					Transform: ballTransform2,
				},
			}

			middleRemoved = true
		}

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)

		for i := 0; i < len(balls); i++ {
			ball := balls[i]
			err = ball.Sprite.Draw(ball.Transform)
			handleError(err)
		}

		//cirnoTransform.Move(geometry.V(400, 400).Scaled(system.DeltaTime()))
		//cirnoSprite.Draw(cirnoTransform)

		batch.Draw()
		layout.Draw()

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

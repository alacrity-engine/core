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
	zBufferNodePool, err := mempool.NewPool[*avltree.AVLNode[float32, render.ZBufferData]](
		func() *avltree.AVLNode[float32, render.ZBufferData] {
			return new(avltree.AVLNode[float32, render.ZBufferData])
		})
	handleError(err)
	zBufferDataPool, err := mempool.NewPool[*collections.AVLTree[int64, *render.Sprite]](
		func() *collections.AVLTree[int64, *render.Sprite] {
			tree, _ := collections.NewAVLTree[int64, *render.Sprite]()
			return tree
		},
	)
	handleError(err)
	zBufferPool, err := mempool.NewPool[*collections.AVLTree[float32, render.ZBufferData]](
		func() *collections.AVLTree[float32, render.ZBufferData] {
			tree, _ := collections.NewAVLTree[float32, render.ZBufferData]()
			return tree
		},
	)
	handleError(err)

	zBufferDataProducer := collections.NewAVLProducer[int64, *render.Sprite](
		zBufferDataPool, zBufferDataNodePool)
	zBufferProducer := collections.NewAVLProducer[float32, render.ZBufferData](
		zBufferPool, zBufferNodePool)

	ballCanvas, err := render.NewCanvas(0, render.Ortho2DStandard(), zBufferProducer, zBufferDataProducer)
	handleError(err)
	err = layout.AddCanvas(ballCanvas)
	handleError(err)

	// Create batch.
	batch, err := render.NewBatch(ballTexture, layout)
	handleError(err)
	balls := make([]*Ball, 0, width*height/
		(imgRGBA.Bounds().Dx()*imgRGBA.Bounds().Dy()))

	// Instantiate all the objects and
	// attach them to the batch.
	/*zCounter := 0

	for i := -float64(width); i < width; i += float64(imgRGBA.Bounds().Dx()) * 2.0 {
		for j := -float64(height); j < height; j += float64(imgRGBA.Bounds().Dy()) * 2.0 {
			ballSprite, err := render.NewSpriteFromTextureAndProgram(
				render.DrawModeStatic, render.DrawModeStatic,
				render.DrawModeStatic, ballTexture, shaderProgram,
				geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
			handleError(err)
			ballSprite.SetZ(float32(zCounter))
			err = ballCanvas.AddSprite(ballSprite)
			handleError(err)
			ballTransform := geometry.NewTransform(nil)
			ballTransform.MoveTo(geometry.V(i, j))

			ball := &Ball{
				Sprite:    ballSprite,
				Transform: ballTransform,
			}
			balls = append(balls, ball)

			err = batch.AttachSprite(ball.Sprite)
			handleError(err)

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
	err = ballCanvas.AddSprite(ballSprite1)
	handleError(err)
	err = ballCanvas.AddSprite(ballSprite2)
	handleError(err)

	ballTransform1 := geometry.NewTransform(nil)
	ballTransform2 := geometry.NewTransform(nil)

	ballSprite1.SetZ(-1)
	ballTransform1.MoveTo(geometry.V(float64(imgRGBA.Bounds().Dx()/2), float64(imgRGBA.Bounds().Dy()/2)))

	ballSprite2.SetZ(1)
	ballTransform2.MoveTo(geometry.V(-float64(imgRGBA.Bounds().Dx()/2), -float64(imgRGBA.Bounds().Dy()/2)))

	err = batch.AttachSprite(ballSprite1)
	handleError(err)
	err = batch.AttachSprite(ballSprite2)
	handleError(err)

	balls = append(balls, &Ball{
		Sprite:    ballSprite1,
		Transform: ballTransform1,
	}, &Ball{
		Sprite:    ballSprite2,
		Transform: ballTransform2,
	})

	system.InitMetrics()

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)

		for i := 0; i < len(balls); i++ {
			ball := balls[i]
			err = ball.Sprite.Draw(ball.Transform)
			handleError(err)
		}

		batch.Draw()

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

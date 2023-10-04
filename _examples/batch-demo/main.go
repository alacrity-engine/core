package main

import (
	"image"
	"image/draw"
	"os"
	"runtime"
	"time"

	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/render"
	"github.com/alacrity-engine/core/system"
	"github.com/alacrity-engine/core/system/collections"
	"github.com/zergon321/go-avltree"
	"github.com/zergon321/mempool"
	"golang.org/x/image/colornames"
)

type Ball struct {
	Sprite     *render.Sprite
	Transform  *geometry.Transform
	MoveVector geometry.Vec
}

const (
	width     = 1920
	height    = 1080
	cooldown  = 200 * time.Millisecond
	dAngle    = 20
	radius    = 60
	numPoints = 10
	speed     = 1
)

func init() {
	runtime.LockOSThread()
}

func main() {
	// Initialize the engine.
	err := system.InitializeWindow("Demo", width, height, true, false)
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
	balls := make(map[*Ball]struct{}, width*height/
		(imgRGBA.Bounds().Dx()*imgRGBA.Bounds().Dy()))

	ticker := time.NewTicker(cooldown)
	points := geometry.PlacingCircle(geometry.ZV, radius, numPoints)

	system.InitMetrics()

	for !system.ShouldClose() {
		system.UpdateDeltaTime()

		if system.ButtonPressed(system.KeyEscape) {
			return
		}

		render.SetClearColor(render.ToRGBA(colornames.Aquamarine))
		render.Clear(render.ClearBitColor | render.ClearBitDepth)

		select {
		case <-ticker.C:
			for i := 0; i < len(points); i++ {
				ballSprite, err := render.NewSpriteFromTextureAndProgram(
					render.DrawModeStatic, render.DrawModeStatic,
					render.DrawModeStatic, ballTexture, shaderProgram,
					geometry.R(0, 0, float64(imgRGBA.Rect.Dx()), float64(imgRGBA.Rect.Dy())))
				ballCanvas.AttachSpriteToBatch(batch, ballSprite)
				handleError(err)

				ballTransform := geometry.NewTransform(nil)
				ballTransform.MoveTo(points[i])

				ball := &Ball{
					Sprite:     ballSprite,
					Transform:  ballTransform,
					MoveVector: points[i].Sub(geometry.ZV),
				}

				balls[ball] = struct{}{}
			}

			for i := 0; i < len(points); i++ {
				points[i] = points[i].RotatedAround(dAngle*geometry.DegToRad, geometry.ZV)
			}

		default:
		}

		ballsToRemove := map[*Ball]struct{}{}

		for ball := range balls {
			if ball.Transform.Position().X < -width/2 ||
				ball.Transform.Position().X > width/2 ||
				ball.Transform.Position().Y < -height/2 ||
				ball.Transform.Position().Y > height/2 {
				ballsToRemove[ball] = struct{}{}

				continue
			}

			ball.Sprite.Draw(ball.Transform)
			ball.Transform.Move(ball.MoveVector.Scaled(speed * system.DeltaTime()))
		}

		batch.Draw()
		layout.Draw()

		for ball := range ballsToRemove {
			delete(balls, ball)
			err = batch.DetachSprite(ball.Sprite)
			handleError(err)
		}

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

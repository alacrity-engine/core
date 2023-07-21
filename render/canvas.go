package render

import (
	"fmt"

	"github.com/alacrity-engine/core/geometry"
	"github.com/alacrity-engine/core/system/collections"
	"github.com/go-gl/mathgl/mgl32"
)

// TODO: implement Z-sorting buffer
// using a balanced binary tree
// (AVL or RB) operating on an array.

type Canvas struct {
	index                   int
	pos                     byte
	sprites                 map[*Sprite]*geometry.Transform
	zBuffer                 collections.SortedDictionary[float32, ZBufferData] // zBuffer is used to draw all the sprites in the order of their Z coordinates.
	layout                  *Layout
	camera                  *Camera
	projection              mgl32.Mat4
	zBufferDataDictProducer collections.SortedDictionaryProducer[int64, *Sprite]
}

type ZBufferData struct {
	lastTimestamp int64
	sprites       map[*Sprite]int64
	timestamps    collections.SortedDictionary[int64, *Sprite]
}

func addSpriteToZBuffer(sprite *Sprite) func(data ZBufferData) (ZBufferData, error) {
	return func(data ZBufferData) (ZBufferData, error) {
		data.lastTimestamp++
		data.sprites[sprite] = data.lastTimestamp
		err := data.timestamps.Add(data.lastTimestamp, sprite)

		if err != nil {
			return ZBufferData{}, err
		}

		return data, nil
	}
}

func removeSpriteFromZBuffer(sprite *Sprite) func(data ZBufferData, found bool) (ZBufferData, error) {
	return func(data ZBufferData, found bool) (ZBufferData, error) {
		if found {
			ts := data.sprites[sprite]
			delete(data.sprites, sprite)
			err := data.timestamps.Remove(ts)

			if err != nil {
				return ZBufferData{}, err
			}
		}

		return data, nil
	}
}

func (canvas *Canvas) newZBufferDataForSprite(sprite *Sprite) (ZBufferData, error) {
	newTree, err := canvas.zBufferDataDictProducer.Produce()

	if err != nil {
		return ZBufferData{}, err
	}

	err = newTree.Add(0, sprite)

	if err != nil {
		return ZBufferData{}, err
	}

	return ZBufferData{
		lastTimestamp: 0,
		sprites: map[*Sprite]int64{
			sprite: 0,
		},
		timestamps: newTree,
	}, nil
}

func (canvas *Canvas) draw() error {
	canvas.zBuffer.VisitInOrder(func(key float32, data ZBufferData) {
		data.timestamps.VisitInOrder(func(key int64, sprite *Sprite) {
			transform := canvas.sprites[sprite]

			if transform != nil {
				sprite.draw(transform.Data(), canvas.
					camera.View(), canvas.projection)

				canvas.sprites[sprite] = nil
			}
		})
	})

	return nil
}

func (canvas *Canvas) Camera() *Camera {
	return canvas.camera
}

func (canvas *Canvas) Index() int {
	return canvas.index
}

func (canvas *Canvas) Z() float32 {
	zLength := zMax - zMin

	return zLength * float32(canvas.index)
}

func (canvas *Canvas) Range() (float32, float32) {
	zLength := zMax - zMin

	return zMin + float32(canvas.index)*zLength,
		zMax + float32(canvas.index)*zLength
}

func (canvas *Canvas) AddSprite(sprite *Sprite) error {
	if _, ok := canvas.sprites[sprite]; ok {
		return fmt.Errorf(
			"the sprite already exists on the canvas")
	}

	canvas.sprites[sprite] = nil
	sprite.canvas = canvas

	// Add the sprite to the Z buffer.
	zData, err := canvas.newZBufferDataForSprite(sprite)

	if err != nil {
		return err
	}

	err = canvas.zBuffer.AddOrUpdate(sprite.drawZ, zData,
		addSpriteToZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) addSpriteFromBatch(sprite *Sprite) error {
	if _, ok := canvas.sprites[sprite]; ok {
		return fmt.Errorf(
			"the sprite already exists on the canvas")
	}

	canvas.sprites[sprite] = nil
	//sprite.canvas = canvas

	// Add the sprite to the Z buffer.
	zData, err := canvas.newZBufferDataForSprite(sprite)

	if err != nil {
		return err
	}

	err = canvas.zBuffer.AddOrUpdate(sprite.drawZ, zData,
		addSpriteToZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) setSpriteZ(sprite *Sprite, oldZ, newZ float32) error {
	err := canvas.zBuffer.Update(oldZ, removeSpriteFromZBuffer(sprite))

	if err != nil {
		return err
	}

	zData, err := canvas.newZBufferDataForSprite(sprite)

	if err != nil {
		return err
	}

	err = canvas.zBuffer.AddOrUpdate(newZ, zData,
		addSpriteToZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) RemoveSprite(sprite *Sprite) error {
	if _, ok := canvas.sprites[sprite]; !ok {
		return fmt.Errorf(
			"the sprite doesn't exist on the canvas")
	}

	delete(canvas.sprites, sprite)
	sprite.canvas = nil

	err := canvas.zBuffer.Update(sprite.drawZ, removeSpriteFromZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) removeBatchedSprite(sprite *Sprite) error {
	if _, ok := canvas.sprites[sprite]; !ok {
		return fmt.Errorf(
			"the sprite doesn't exist on the canvas")
	}

	delete(canvas.sprites, sprite)
	//sprite.canvas = nil

	err := canvas.zBuffer.Update(sprite.drawZ, removeSpriteFromZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) updateBatchViews() {
	for i := 0; i < len(canvas.layout.batches); i++ {
		batch := canvas.layout.batches[i]
		batch.setCanvasView(int(canvas.pos), canvas.camera.View())
	}
}

func NewCanvas(
	drawZ int, projection mgl32.Mat4,
	zBufferDictProducer collections.SortedDictionaryProducer[float32, ZBufferData],
	zBufferDataDictProducer collections.SortedDictionaryProducer[int64, *Sprite],
) (*Canvas, error) {
	camera := NewCamera()
	canvas := &Canvas{
		sprites:                 map[*Sprite]*geometry.Transform{},
		index:                   drawZ,
		camera:                  NewCamera(),
		projection:              projection,
		zBufferDataDictProducer: zBufferDataDictProducer,
	}

	camera.canvas = canvas
	zBuffer, err := zBufferDictProducer.Produce()

	if err != nil {
		return nil, err
	}

	canvas.zBuffer = zBuffer

	return canvas, nil
}

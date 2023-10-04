package render

import (
	"fmt"

	"github.com/alacrity-engine/core/math/geometry"
	"github.com/alacrity-engine/core/system/collections"
	"github.com/go-gl/mathgl/mgl32"
)

type Canvas struct {
	index                   int
	pos                     byte
	name                    string
	sprites                 map[*Sprite]*geometry.Transform
	batches                 map[*Batch]bool
	batchNameIndex          map[string]*Batch
	zBuffer                 collections.UnrestrictedSortedDictionary[Geometric, ZBufferData] // zBuffer is used to draw all the sprites in the order of their Z coordinates.
	layout                  *Layout
	camera                  *Camera
	projection              mgl32.Mat4
	zBufferDataDictProducer collections.SortedDictionaryProducer[int64, *Sprite]
}

type ZBufferData struct {
	lastTimestamp int64
	sprites       map[*Sprite]int64
	timestamps    collections.SortedDictionary[int64, *Sprite]
	batch         *Batch
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
	err := canvas.zBuffer.VisitInOrder(func(key Geometric, data ZBufferData) error {
		if len(data.sprites) > 0 {
			data.timestamps.VisitInOrder(func(key int64, sprite *Sprite) error {
				transform := canvas.sprites[sprite]

				if transform != nil {
					sprite.draw(transform.Data(), canvas.
						camera.View(), canvas.projection)

					canvas.sprites[sprite] = nil
				}

				return nil
			})
		} else if data.batch != nil && canvas.batches[data.batch] {
			data.batch.draw()
			canvas.batches[data.batch] = false
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) Name() string {
	return canvas.name
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

	err = canvas.zBuffer.AddOrUpdate(Point{Z: sprite.drawZ}, zData,
		addSpriteToZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) AddBatch(batch *Batch, z1, z2 float32) error {
	if _, ok := canvas.batchNameIndex[batch.name]; ok {
		return fmt.Errorf(
			"a batch named '%s' already exists on the canvas", batch.name)
	}

	if z1 < zMin {
		return fmt.Errorf("the Z1=%f is less than %f", z1, zMin)
	}

	if z2 > zMax {
		return fmt.Errorf("the Z2=%f is less than %f", z2, zMax)
	}

	if _, ok := canvas.batches[batch]; ok {
		return fmt.Errorf("the batch already exists on the canvas")
	}

	batch.canvas = canvas
	canvas.batches[batch] = false
	batch.z1 = z1
	batch.z2 = z2

	data := ZBufferData{
		batch: batch,
	}

	err := canvas.zBuffer.AddOrUpdate(Range{Z1: z1, Z2: z2}, data,
		func(oldValue ZBufferData) (ZBufferData, error) {
			return ZBufferData{}, fmt.Errorf(
				"the batch intersects with existing objects")
		})

	if err != nil {
		return err
	}

	canvas.batchNameIndex[batch.name] = batch

	return nil
}

func (canvas *Canvas) AttachSpriteToBatch(batch *Batch, sprite *Sprite) error {
	sprite.canvas = canvas
	return batch.attachSprite(sprite)
}

func (canvas *Canvas) setSpriteZ(sprite *Sprite, oldZ, newZ float32) error {
	err := canvas.zBuffer.Update(Point{Z: oldZ}, removeSpriteFromZBuffer(sprite))

	if err != nil {
		return err
	}

	zData, err := canvas.newZBufferDataForSprite(sprite)

	if err != nil {
		return err
	}

	err = canvas.zBuffer.AddOrUpdate(Point{Z: newZ}, zData,
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

	err := canvas.zBuffer.Update(Point{Z: sprite.drawZ},
		removeSpriteFromZBuffer(sprite))

	if err != nil {
		return err
	}

	return nil
}

func (canvas *Canvas) RemoveBatch(batch *Batch) error {
	if _, ok := canvas.batches[batch]; !ok {
		return fmt.Errorf("the batch doesn't exist on the canvas")
	}

	z1 := batch.z1
	z2 := batch.z2

	batch.canvas = nil
	delete(canvas.batches, batch)
	batch.z1 = 0
	batch.z2 = 0

	err := canvas.zBuffer.Remove(Range{Z1: z1, Z2: z2})

	if err != nil {
		return err
	}

	return nil
}

func NewCanvas(
	name string, drawZ int, projection mgl32.Mat4,
	zBufferDictProducer collections.UnrestrictedSortedDictionaryProducer[Geometric, ZBufferData],
	zBufferDataDictProducer collections.SortedDictionaryProducer[int64, *Sprite],
) (*Canvas, error) {
	camera := NewCamera()
	canvas := &Canvas{
		name:                    name,
		sprites:                 map[*Sprite]*geometry.Transform{},
		batches:                 map[*Batch]bool{},
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

package resources

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/alacrity-engine/core/anim"
	"github.com/alacrity-engine/core/render"

	codec "github.com/alacrity-engine/resource-codec"
	bolt "go.etcd.io/bbolt"
)

// TODO: load and cache shader sources,
// textures, pre-compiled shaders.

// TODO: animations must be based
// on textures, not pictures. Anytime
// the the engine requests an animation
// by id, the resource loader should
// create a new instance of the requested
// animation.

// ResourceLoader loads sprites,
// animations, sound and text
// from resource files.
type ResourceLoader struct {
	resourceFile *bolt.DB
	buffer       *resourceBuffer
}

// Close closes the resource file.
func (loader *ResourceLoader) Close() error {
	return loader.resourceFile.Close()
}

// LoadAnimation loads the animation with spritesheet
// and frames from the resource file.
func (loader *ResourceLoader) LoadAnimation(animID string) (*anim.Animation, error) {
	// Load the animation frames from the buffer
	// or the resource file.
	animData, err := loader.buffer.takeAnimation(animID)

	if err != nil {
		switch err.(type) {
		case *ErrorAnimationDoesntExist:
			err = loader.resourceFile.View(func(tx *bolt.Tx) error {
				buck := tx.Bucket([]byte("animations"))

				if buck == nil {
					return fmt.Errorf(
						"the 'animations' bucket doesn't exist")
				}

				data := buck.Get([]byte(animID))

				if data == nil {
					return fmt.Errorf("no '%s' animation", animID)
				}

				animData, err = codec.AnimationDataFromBytes(data)

				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return nil, err
			}

			err = loader.buffer.putAnimation(animID, animData)

			if err != nil {
				return nil, err
			}

		default:
			return nil, err
		}
	}

	texture, err := loader.LoadTexture(animData.TextureID)

	if err != nil {
		return nil, err
	}

	delays := []time.Duration{}

	for _, duration := range animData.Durations {
		delay := time.Duration(duration) * time.Millisecond
		delays = append(delays, delay)
	}

	anim, err := anim.NewAnimation(
		texture, animData.Frames, delays, false)

	if err != nil {
		return nil, err
	}

	return anim, nil
}

func (loader *ResourceLoader) LoadTexture(name string) (*render.Texture, error) {
	texture, err := loader.buffer.takeTexture(name)

	if err != nil {
		switch err.(type) {
		case *ErrorTextureDoesntExist:
			var texData *codec.TextureData

			err = loader.resourceFile.View(func(tx *bolt.Tx) error {
				buck := tx.Bucket([]byte("textures"))

				if buck == nil {
					return fmt.Errorf("bucket 'textures' not found")
				}

				data := buck.Get([]byte(name))

				if data == nil {
					return fmt.Errorf("no '%s' texture", name)
				}

				texData, err = codec.TextureDataFromBytes(data)

				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return nil, err
			}

			pic, err := loader.LoadPicture(texData.PictureID)

			if err != nil {
				return nil, err
			}

			texture = render.NewTextureFromPicture(
				pic, render.TextureFiltering(texData.Filtering))
			err = loader.buffer.putTexture(name, texture)

			if err != nil {
				return nil, err
			}

		default:
			return nil, err
		}
	}

	return texture, nil
}

func (loader *ResourceLoader) loadSpritesheet(id string) (*codec.SpritesheetData, error) {
	ss, err := loader.buffer.takeSpritesheet(id)

	if err != nil {
		switch err.(type) {
		case *ErrorSpritesheetDoesntExist:
			err = loader.resourceFile.View(func(tx *bolt.Tx) error {
				buck := tx.Bucket([]byte("spritesheets"))

				if buck == nil {
					return fmt.Errorf("no 'spritesheets' bucket found")
				}

				data := buck.Get([]byte(id))

				if data == nil {
					return fmt.Errorf("no '%s' spritesheet found", id)
				}

				ss, err = codec.SpritesheetDataFromBytes(data)

				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return nil, err
			}

			err = loader.buffer.putSpritesheet(id, ss)

			if err != nil {
				return nil, err
			}

		default:
			return nil, err
		}
	}

	return ss, nil
}

// LoadPicture loads the picture from the resource file by the name of the picture.
func (loader *ResourceLoader) LoadPicture(name string) (*render.Picture, error) {
	picture, err := loader.buffer.takePicture(name)

	if err != nil {
		switch err.(type) {
		case *ErrorPictureDoesntExist:
			er := loader.resourceFile.View(func(tx *bolt.Tx) error {
				buck := tx.Bucket([]byte("spritesheets"))

				if buck == nil {
					return fmt.Errorf("bucket 'spritesheets' not found")
				}

				pictureBytes := buck.Get([]byte(name))

				if pictureBytes == nil {
					return fmt.Errorf("picture with ID '%s' not found",
						name)
				}

				compressedPicture, err := codec.CompressedPictureFromBytes(pictureBytes)

				if err != nil {
					return err
				}

				pictureDataDec, err := compressedPicture.Decompress()

				if err != nil {
					return err
				}

				picture = PictureDataToPicture(pictureDataDec)

				return nil
			})

			if er != nil {
				return nil, er
			}

			er = loader.buffer.putPicture(name, picture)

			if er != nil {
				return nil, er
			}

		default:
			return nil, err
		}
	}

	return picture, nil
}

// LoadFont loads a font stored in the resource file under the specified name.
// func (loader *ResourceLoader) LoadFont(name string) (*truetype.Font, error) {
// 	font, err := loader.buffer.takeFont(name)
//
// 	if err != nil {
// 		switch err.(type) {
// 		case *ErrorFontDoesntExist:
// 			er := loader.resourceFile.View(func(tx *bolt.Tx) error {
// 				buck := tx.Bucket([]byte("fonts"))
//
// 				if buck == nil {
// 					return fmt.Errorf("bucket 'fonts' not found")
// 				}
//
// 				fontData := buck.Get([]byte(name))
//
// 				if fontData == nil {
// 					return fmt.Errorf("font '%s' not found", name)
// 				}
//
// 				var err error
// 				font, err = truetype.Parse(fontData)
//
// 				if err != nil {
// 					return err
// 				}
//
// 				return nil
// 			})
//
// 			if er != nil {
// 				return nil, er
// 			}
//
// 			er = loader.buffer.putFont(name, font)
//
// 			if er != nil {
// 				return nil, er
// 			}
//
// 		default:
// 			return nil, err
// 		}
// 	}
//
// 	return font, nil
// }

// LoadAudio loads the specified audio from the resource file.
func (loader *ResourceLoader) LoadAudio(name string) (io.ReadCloser, error) {
	audio, err := loader.buffer.takeAudio(name)

	if err != nil {
		switch err.(type) {
		case *ErrorAudioDoesntExist:
			er := loader.resourceFile.View(func(tx *bolt.Tx) error {
				buck := tx.Bucket([]byte("audio"))

				if buck == nil {
					return fmt.Errorf("bucket 'audio' not found")
				}

				audio = buck.Get([]byte(name))

				if audio == nil {
					return fmt.Errorf("audio '%s' not found", name)
				}

				return nil
			})

			if er != nil {
				return nil, er
			}

			er = loader.buffer.putAudio(name, audio)

			if er != nil {
				return nil, er
			}

		default:
			return nil, err
		}
	}

	reader := bytes.NewReader(audio)
	stream := io.NopCloser(reader)

	return stream, nil
}

// NewResourceLoader crates a new resource loader for the specified resource file.
func NewResourceLoader(file string) (*ResourceLoader, error) {
	resourceFile, err := bolt.Open(file, 0666, nil)

	if err != nil {
		return nil, err
	}

	return &ResourceLoader{
		resourceFile: resourceFile,
		buffer:       newResourceBuffer(),
	}, nil
}

package resources

import (
	"github.com/alacrity-engine/core/render"
	codec "github.com/alacrity-engine/resource-codec"
)

func PictureDataToPicture(picData *codec.PictureData) *render.Picture {
	picture := &render.Picture{
		Width:  picData.Width,
		Height: picData.Height,
		Pix:    make([]byte, len(picData.Pix)),
	}

	copy(picture.Pix, picData.Pix)

	return picture
}

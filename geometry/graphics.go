package geometry

func ComputeSpriteVertices(width, height int, targetArea Rect) []float32 {
	texToScreenWidth := float32(targetArea.W() / float64(width))
	texToscreenHeight := float32(targetArea.H() / float64(width))
	vertices := []float32{
		texToScreenWidth * -1.0, texToscreenHeight * -1.0, 0.0,
		texToScreenWidth * -1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0,
	}

	return vertices
}

func ComputeSpriteVerticesFill(buffer []float32, width, height int, targetArea Rect) {
	texToScreenWidth := float32(targetArea.W() / float64(width))
	texToscreenHeight := float32(targetArea.H() / float64(width))
	vertices := []float32{
		texToScreenWidth * -1.0, texToscreenHeight * -1.0, 0.0,
		texToScreenWidth * -1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0,
	}

	copy(buffer, vertices)
}

func ComputeSpriteTextureCoordinates(imageWidth, imageHeight int, targetArea Rect) []float32 {
	return []float32{
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Min.Y) / float32(imageHeight),
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight),
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight),
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Min.Y) / float32(imageHeight),
	}
}

func ComputeSpriteTextureCoordinatesFill(buffer []float32, imageWidth, imageHeight int, targetArea Rect) {
	texCoords := []float32{
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Min.Y) / float32(imageHeight),
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight),
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight),
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Min.Y) / float32(imageHeight),
	}

	copy(buffer, texCoords)
}

func ComputeSpriteVerticesNoElementsFill(buffer []float32, width, height int, targetArea Rect) {
	texToScreenWidth := float32(targetArea.W() / float64(width))
	texToscreenHeight := float32(targetArea.H() / float64(width))
	vertices := []float32{
		texToScreenWidth * -1.0, texToscreenHeight * -1.0, 0.0,
		texToScreenWidth * -1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0, // extraneous 3
		texToScreenWidth * -1.0, texToscreenHeight * 1.0, 0.0, // extraneous 1
		texToScreenWidth * 1.0, texToscreenHeight * 1.0, 0.0,
		texToScreenWidth * 1.0, texToscreenHeight * -1.0, 0.0,
	}

	copy(buffer, vertices)
}

func ComputeSpriteTextureCoordinatesNoElementsFill(buffer []float32, imageWidth, imageHeight int, targetArea Rect) {
	coords := []float32{
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Min.Y) / float32(imageHeight),
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight),
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight), // extraneous 3
		float32(targetArea.Min.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight), // extraneous 1
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Max.Y) / float32(imageHeight),
		float32(targetArea.Max.X) / float32(imageWidth), float32(targetArea.Min.Y) / float32(imageHeight),
	}

	copy(buffer, coords)
}

func ColorMaskDataNoElementsFill(buffer []float32, colorMaskData [16]float32) {
	copy(buffer, colorMaskData[0:8])
	copy(buffer[8:12], colorMaskData[12:16]) // extraneous 3
	copy(buffer[12:16], colorMaskData[4:8])  // extraneous 1
	copy(buffer[16:24], colorMaskData[8:16])
}

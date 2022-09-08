package render

import "fmt"

var (
	width  int
	height int
)

func Width() int {
	return width
}

func Height() int {
	return height
}

func SetWidth(_width int) error {
	if _width <= 0 {
		return fmt.Errorf("width must be above 0")
	}

	width = _width

	return nil
}

func SetHeight(_height int) error {
	if _height <= 0 {
		return fmt.Errorf("height must be above 0")
	}

	height = _height

	return nil
}

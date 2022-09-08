package debug

import (
	"fmt"
	"image/color"

	"github.com/alacrity-engine/core/resources"
	"github.com/alacrity-engine/core/system"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	colors "golang.org/x/image/colornames"
)

var (
	txt *text.Text
)

// ConsoleInitialize initializes the debug console.
func ConsoleInitialize() error {
	// Load a new font.
	resourceLoader, err := resources.NewResourceLoader("common.res")

	if err != nil {
		return err
	}
	defer resourceLoader.Close()

	fnt, err := resourceLoader.LoadFont("debug")

	if err != nil {
		return err
	}

	face := truetype.NewFace(fnt, &truetype.Options{
		Size:              20,
		GlyphCacheEntries: 1,
	})

	// Initialize the debug console.
	atlas := text.NewAtlas(face, text.ASCII)
	txt = text.New(pixel.V(20, 1040), atlas)

	txt.Color = colors.White

	return nil
}

// ConsoleClear clears the debug console.
func ConsoleClear() {
	if txt == nil {
		return
	}

	txt.Clear()
}

// ConsolePrintln prints a new line to the debug console.
func ConsolePrintln(args ...interface{}) {
	if txt == nil {
		return
	}

	fmt.Fprintln(txt, args...)
}

// ConsolePrintf performs a formatted print to the debug console.
func ConsolePrintf(message string, args ...interface{}) {
	if txt == nil {
		return
	}

	fmt.Fprintf(txt, message, args...)
}

// ConsolePrint prints the arguments in the console.
func ConsolePrint(args ...interface{}) {
	if txt == nil {
		return
	}

	fmt.Fprint(txt, args...)
}

// ConsoleSetColor sets the text color
// in the debug console.
func ConsoleSetColor(cl color.RGBA) {
	if txt == nil {
		return
	}

	txt.Color = cl
}

// ConsoleColor returns the text color in the console.
func ConsoleColor() color.RGBA {
	if txt == nil {
		return colors.White
	}

	return txt.Color.(color.RGBA)
}

// ConsoleOutput outputs the console contents
// to the game window.
func ConsoleOutput() {
	if txt == nil {
		return
	}

	txt.Draw(system.Window(), pixel.IM)
}

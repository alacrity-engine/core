package system

// #cgo LDFLAGS: -lX11 -lXrandr -lstdc++
// #include <X11/Xlib.h>
// #include <X11/extensions/Xrandr.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// DisplayResolution returns the current width
// and height of the main display.
func DisplayResolution() (int, int, error) {
	// Get the main display of X server.
	display := C.XOpenDisplay(nil)
	// Get the ID of the root window.
	window := C.XRootWindow(display, 0)
	// Get the information about the display.
	conf := C.XRRGetScreenInfo(display, window)
	// Get the array of resolutions.
	var numSizes C.int
	resolutions := C.XRRSizes(display, 0, &numSizes)
	// Turn it into a Go slice.
	goResolutions := (*[1 << 30]C.XRRScreenSize)(
		unsafe.Pointer(resolutions))[:numSizes:numSizes]
	// Get the rotation of the display.
	var rotation C.Rotation
	// Get the ID of the display resolution in the table of resolutions.
	resolutionIdx := C.XRRConfigCurrentConfiguration(conf, &rotation)
	// Get the resolution by the ID.
	resolution := goResolutions[resolutionIdx]

	C.XCloseDisplay(display)

	return int(resolution.width), int(resolution.height), nil
}

// SetDisplayResolution sets the resolution of
// the display to the required width and height.
func SetDisplayResolution(width, height int) error {
	// Get the main display of X server.
	display := C.XOpenDisplay(nil)
	// Get the ID of the root window.
	window := C.XRootWindow(display, 0)
	// Get the information about the display.
	conf := C.XRRGetScreenInfo(display, window)
	// Get the rotation of the display.
	var originalRotation C.Rotation
	C.XRRConfigCurrentConfiguration(conf, &originalRotation)
	// Get the array of resolutions.
	var numSizes C.int
	resolutions := C.XRRSizes(display, 0, &numSizes)
	// Turn it into a Go slice.
	goResolutions := (*[1 << 30]C.XRRScreenSize)(
		unsafe.Pointer(resolutions))[:numSizes:numSizes]

	// Get all the frequency rates for each resolution.
	resRates := make([][]C.short, numSizes)

	for i := 0; i < int(numSizes); i++ {
		var numRates C.int
		rates := C.XRRRates(display, 0, C.int(i), &numRates)
		goRates := (*[1 << 30]C.short)(
			unsafe.Pointer(rates))[:numSizes:numSizes]

		// Add the array of rates to the table
		// of rates for resolutions.
		resRates[i] = goRates
	}

	// Check if the requested resolution
	// is supported by the display.
	resInd := -1

	for i := 0; i < int(numSizes); i++ {
		if int(goResolutions[i].width) == width &&
			int(goResolutions[i].height) == height {
			resInd = i
			break
		}
	}

	if resInd < 0 {
		return fmt.Errorf("resolution %dx%d is not supported by the display",
			width, height)
	}

	C.XRRSetScreenConfigAndRate(display, conf, window, C.int(resInd),
		originalRotation, resRates[resInd][0], C.CurrentTime)
	C.XCloseDisplay(display)

	return nil
}

package system

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	// deviceNameLength is the length of the
	// device name in DevMode structure.
	deviceNameLength = 32
	// formNameLength is the length of the
	// form name in DevMode structure.
	formNameLength = 32
	// enumCurrentSettings is a command to enumerate
	// the current display settings.
	enumCurrentSettings uint32 = 0xFFFFFFFF
	// enumCurrentSettings is a command to enumerate
	// the display settings saved in the system registry.
	enumRegistrySettings uint32 = 0xFFFFFFFE
	// displayChangeSuccessful status code indicates
	// that the display change attempt was successful.
	displayChangeSuccessful uint32 = 0
	// displayChangeRestart indicates that the computer needs
	// to be restarted in order to change the display settings.
	displayChangeRestart uint32 = 1
	// displayChangeFailed indicates that the system
	// failed to change the display settings.
	displayChangeFailed uint32 = 0xFFFFFFFF
	// displayChangeBadMode indicates that the
	// requested resolution is not supported.
	displayChangeBadMode uint32 = 0xFFFFFFFE
)

// DevMode is a structure used to
// specify characteristics of display
// and print devices.
type DevMode struct {
	DeviceName       [deviceNameLength]uint16
	SpecVersion      uint16
	DriverVersion    uint16
	Size             uint16
	DriverExtra      uint16
	Fields           uint32
	Orientation      int16
	PaperSize        int16
	PaperLength      int16
	PaperWidth       int16
	Scale            int16
	Copies           int16
	DefaultSource    int16
	PrintQuality     int16
	Color            int16
	Duplex           int16
	YResolution      int16
	TTOption         int16
	Collate          int16
	FormName         [formNameLength]uint16
	LogPixels        uint16
	BitsPerPel       uint32
	PelsWidth        uint32
	PelsHeight       uint32
	DisplayFlags     uint32
	DisplayFrequency uint32
	ICMMethod        uint32
	ICMIntent        uint32
	MediaType        uint32
	DitherType       uint32
	Reserved1        uint32
	Reserved2        uint32
	PanningWidth     uint32
	PanningHeight    uint32
}

// DisplayResolution returns the current width
// and height of the main display.
func DisplayResolution() (int, int, error) {
	// Load the DLL and the procedures.
	user32dll := syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW := user32dll.NewProc("EnumDisplaySettingsW")

	// Get the display information.
	devMode := new(DevMode)
	ret, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(nil)),
		uintptr(enumCurrentSettings), uintptr(unsafe.Pointer(devMode)))

	if ret == 0 {
		return -1, -1, fmt.Errorf("couldn't get the screen resolution")
	}

	return int(devMode.PelsWidth), int(devMode.PelsHeight), nil
}

// SetDisplayResolution sets the resolution of
// the display to the required width and height.
func SetDisplayResolution(width, height int) error {
	// Load the DLL and the procedures.
	user32dll := syscall.NewLazyDLL("user32.dll")
	procEnumDisplaySettingsW := user32dll.NewProc("EnumDisplaySettingsW")
	procChangeDisplaySettingsW := user32dll.NewProc("ChangeDisplaySettingsW")

	// Get the display information.
	devMode := new(DevMode)
	statusCode, _, _ := procEnumDisplaySettingsW.Call(uintptr(unsafe.Pointer(nil)),
		uintptr(enumCurrentSettings), uintptr(unsafe.Pointer(devMode)))

	if statusCode == 0 {
		return fmt.Errorf("couldn't get the screen resolution")
	}

	newMode := *devMode
	newMode.PelsWidth = uint32(width)
	newMode.PelsHeight = uint32(height)
	statusCode, _, _ = procChangeDisplaySettingsW.Call(uintptr(unsafe.Pointer(&newMode)),
		uintptr(0))

	switch statusCode {
	case uintptr(displayChangeSuccessful):
		return nil

	case uintptr(displayChangeRestart):
		return fmt.Errorf("system reboot required to change the screen resolution")

	case uintptr(displayChangeBadMode):
		return fmt.Errorf("the resolution is not supported by the display")

	case uintptr(displayChangeFailed):
		return fmt.Errorf("failed to change the display resolution")
	}

	return nil
}

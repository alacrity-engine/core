package resources

import "fmt"

// ErrorPictureAlreadyExists is raised when
// the resource loader tries to put the picture
// that is already loaded in the buffer.
type ErrorPictureAlreadyExists struct {
	picName string
}

// Error returns the error message.
func (err *ErrorPictureAlreadyExists) Error() string {
	return fmt.Sprintf("picture with name %s already exists in the resource buffer",
		err.picName)
}

// RaiseErrorPictureAlreadyExists returns a new error
// about picture that is already loaded.
func RaiseErrorPictureAlreadyExists(picName string) *ErrorPictureAlreadyExists {
	return &ErrorPictureAlreadyExists{
		picName: picName,
	}
}

/*****************************************************************************************************************/

// ErrorPictureDoesntExist is raised when
// the resource loader tries to take the picture
// that isn't yet loaded.
type ErrorPictureDoesntExist struct {
	picName string
}

// Error returns the error message.
func (err *ErrorPictureDoesntExist) Error() string {
	return fmt.Sprintf("picture with name %s doesn't exist in the resource buffer",
		err.picName)
}

// RaiseErrorPictureDoesntExist returns a new error
// about the picture that is not yet loaded.
func RaiseErrorPictureDoesntExist(picName string) *ErrorPictureDoesntExist {
	return &ErrorPictureDoesntExist{
		picName: picName,
	}
}

/*****************************************************************************************************************/

// ErrorAnimationAlreadyExists is raised when
// the resource loader tries to put the animation
// that is already loaded in the buffer.
type ErrorAnimationAlreadyExists struct {
	animName string
}

// Error returns the error message.
func (err *ErrorAnimationAlreadyExists) Error() string {
	return fmt.Sprintf("animation with name %s already exists in the resource buffer",
		err.animName)
}

// RaiseErrorAnimationAlreadyExists returns a new error
// about the animation that is already loaded.
func RaiseErrorAnimationAlreadyExists(animName string) *ErrorAnimationAlreadyExists {
	return &ErrorAnimationAlreadyExists{
		animName: animName,
	}
}

/*****************************************************************************************************************/

// ErrorAnimationDoesntExist is raised when
// the resource loader tries to take the animation
// that isn't yet loaded.
type ErrorAnimationDoesntExist struct {
	animName string
}

// Error returns the error message.
func (err *ErrorAnimationDoesntExist) Error() string {
	return fmt.Sprintf("animation with name %s doesn't exist in the resource buffer",
		err.animName)
}

// RaiseErrorAnimationDoesntExist returns a new error
// about the animation that is not yet loaded.
func RaiseErrorAnimationDoesntExist(animName string) *ErrorAnimationDoesntExist {
	return &ErrorAnimationDoesntExist{
		animName: animName,
	}
}

/*****************************************************************************************************************/

// ErrorFontAlreadyExists is raised when
// the resource loader tries to put the font
// that is already loaded in the buffer.
type ErrorFontAlreadyExists struct {
	fontName string
}

// Error returns the error message.
func (err *ErrorFontAlreadyExists) Error() string {
	return fmt.Sprintf("font with name %s already exists in the resource buffer",
		err.fontName)
}

// RaiseErrorFontAlreadyExists returns a new error
// about the font that is already loaded.
func RaiseErrorFontAlreadyExists(fontName string) *ErrorFontAlreadyExists {
	return &ErrorFontAlreadyExists{
		fontName: fontName,
	}
}

/*****************************************************************************************************************/

// ErrorFontDoesntExist is raised when
// the resource loader tries to take the font
// that isn't yet loaded.
type ErrorFontDoesntExist struct {
	picName string
}

// Error returns the error message.
func (err *ErrorFontDoesntExist) Error() string {
	return fmt.Sprintf("font with name %s doesn't exist in the resource buffer",
		err.picName)
}

// RaiseErrorFontDoesntExist returns a new error
// about the font that is not yet loaded.
func RaiseErrorFontDoesntExist(picName string) *ErrorFontDoesntExist {
	return &ErrorFontDoesntExist{
		picName: picName,
	}
}

/*****************************************************************************************************************/

// ErrorAudioAlreadyExists is raised when
// the resource loader tries to put the audio
// that is already loaded in the buffer.
type ErrorAudioAlreadyExists struct {
	audioName string
}

// Error returns the error message.
func (err *ErrorAudioAlreadyExists) Error() string {
	return fmt.Sprintf("audio with name %s already exists in the resource buffer",
		err.audioName)
}

// RaiseErrorAudioAlreadyExists returns a new error
// about the audio that is already loaded.
func RaiseErrorAudioAlreadyExists(audioName string) *ErrorAudioAlreadyExists {
	return &ErrorAudioAlreadyExists{
		audioName: audioName,
	}
}

/*****************************************************************************************************************/

// ErrorAudioDoesntExist is raised when
// the resource loader tries to extract
// the audio from the resource buffer and
// it doesn't exist.
type ErrorAudioDoesntExist struct {
	audioName string
}

// Error returns the error message.
func (err *ErrorAudioDoesntExist) Error() string {
	return fmt.Sprintf("audio with name %s doesn't exist in the resource buffer",
		err.audioName)
}

// RaiseErrorAudioAlreadyExists returns a new error
// about the audio that doesn't exist in the resource buffer.
func RaiseErrorAudioDoesntExist(audioName string) *ErrorAudioDoesntExist {
	return &ErrorAudioDoesntExist{
		audioName: audioName,
	}
}

/*****************************************************************************************************************/

type ErrorTextureDoesntExist struct {
	textureName string
}

func (err *ErrorTextureDoesntExist) Error() string {
	return fmt.Sprintf("the '%s' texture doesn't exist", err.textureName)
}

func RaiseErrorTextureDoesntExist(textureName string) *ErrorTextureDoesntExist {
	return &ErrorTextureDoesntExist{
		textureName: textureName,
	}
}

/*****************************************************************************************************************/

type ErrorSpritesheetDoesntExist struct {
	ssName string
}

func (err *ErrorSpritesheetDoesntExist) Error() string {
	return fmt.Sprintf("the '%s' spritesheet doesn't exist", err.ssName)
}

func RaiseErrorSpritesheetDoesntExist(ssName string) *ErrorSpritesheetDoesntExist {
	return &ErrorSpritesheetDoesntExist{
		ssName: ssName,
	}
}

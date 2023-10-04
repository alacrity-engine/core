package audio

import (
	"fmt"
	"time"
)

// ErrorVolumeNegative is raised when
// the amount of points to adjust
// the sound is negative.
type ErrorVolumeNegative struct {
	value int
}

// Error returns the error message.
func (err *ErrorVolumeNegative) Error() string {
	return fmt.Sprintf(
		"the volume value is below 0: %d", err.value)
}

// RaiseErrorVolumeNegative returns a new error
// about negative amount of points to change the volume.
func RaiseErrorVolumeNegative(value int) *ErrorVolumeNegative {
	return &ErrorVolumeNegative{
		value: value,
	}
}

/*****************************************************************************************************************/

// ErrorWrongDuration is raised when
// the duration to rewind is higher
// than the actual length of the audio stream.
type ErrorWrongDuration struct {
	value time.Duration
}

// Error returns the error message.
func (err *ErrorWrongDuration) Error() string {
	return fmt.Sprintf(
		"cannot rewind to the duration: %v", err.value)
}

// RaiseErrorWrongDuration raises a new error
// about bad rewind duration.
func RaiseErrorWrongDuration(value time.Duration) *ErrorWrongDuration {
	return &ErrorWrongDuration{
		value: value,
	}
}

package system

import (
	"time"
)

var (
	frameCount int
	perSecond  <-chan time.Time
	fps        int
	lastFrame  time.Time
	deltaTime  float64
)

// InitMetrics initializes metric variables.
func InitMetrics() {
	lastFrame = time.Now()
	fps = 0
	perSecond = time.Tick(time.Second)
}

// UpdateDeltaTime sets new value to the dt variable.
// This method should be called in the start of the frame.
func UpdateDeltaTime() {
	deltaTime = time.Since(lastFrame).Seconds()
	lastFrame = time.Now()
}

// UpdateFrameRate increments the
// frame counter and updates the FPS variable.
// This method should be called in the end of the frame.
func UpdateFrameRate() {
	frameCount++

	select {
	case <-perSecond:
		fps = frameCount
		frameCount = 0

	default:
	}
}

// DeltaTime is the time (in seconds) passed
// since the last frame.
func DeltaTime() float64 {
	return deltaTime
}

// FPS is the number of frames
// per second.
func FPS() int {
	return fps
}

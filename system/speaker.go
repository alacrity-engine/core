package system

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

const (
	// SpeakerSampleRate is the default sample
	// rate of the system speaker.
	SpeakerSampleRate beep.SampleRate = 44100
)

// SpeakerInit initializes the speaker
// with the default sample rate of 44100.
func SpeakerInit() error {
	return speaker.Init(SpeakerSampleRate,
		SpeakerSampleRate.N(time.Second/10))
}

// SpeakerLock locks the
// system speaker.
func SpeakerLock() {
	speaker.Lock()
}

// SpeakerUnlock unlocks the
// system speaker.
func SpeakerUnlock() {
	speaker.Unlock()
}

// SpeakerPlay makes the system speaker
// play all the provided audio streams
// in parallel manner.
func SpeakerPlay(streamers ...beep.Streamer) {
	speaker.Play(streamers...)
}

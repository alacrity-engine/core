package audio

import (
	"io"
	"sync"
	"time"

	"github.com/alacrity-engine/core/engine"
	"github.com/alacrity-engine/core/system"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
)

// AudioSource is a component for
// game object to play sound.
type AudioSource struct {
	engine.BaseComponent
	format            beep.Format
	streamer          beep.StreamSeekCloser
	resampledStreamer *beep.Resampler
	control           *beep.Ctrl
	volumeControl     *effects.Volume
	loop              bool
	loopLocker        *sync.Mutex
	volumeLevel       int
	loopDone          chan bool
	loopCancel        chan bool
}

// Loop returns true if the current
// audio stream should be repeated.
func (as *AudioSource) Loop() bool {
	as.loopLocker.Lock()
	defer as.loopLocker.Unlock()

	return as.loop
}

// SetLoop sets the current audio
// stream to be repeated or not.
func (as *AudioSource) SetLoop(loop bool) {
	as.loopLocker.Lock()
	defer as.loopLocker.Unlock()

	as.loop = loop
}

// Pause pauses the audio playback.
func (as *AudioSource) Pause() {
	system.SpeakerLock()
	as.control.Paused = true
	system.SpeakerUnlock()
}

// Unpause unpauses the audio playback.
func (as *AudioSource) Unpause() {
	system.SpeakerLock()
	as.control.Paused = false
	system.SpeakerUnlock()
}

// VolumeUp increases the volume level by
// the specified amount of points.
//
// The max level of volume is 100, the min
// level of volume is 0. Each point increases
// the volume level by 1%.
func (as *AudioSource) VolumeUp(amount int) error {
	if amount < 0 {
		return RaiseErrorVolumeNegative(amount)
	}

	if amount > 100 || as.volumeLevel+amount > 100 {
		amount = 100 - as.volumeLevel
	}

	if amount != 0 {
		system.SpeakerLock()
		as.volumeControl.Volume += float64(amount)
		system.SpeakerUnlock()
	}

	return nil
}

// VolumeDown decreases the volume level by
// the specified amount of points.
//
// The max level of volume is 100, the min
// level of volume is 0. Each point increases
// the volume level by 1%.
func (as *AudioSource) VolumeDown(amount int) error {
	if amount < 0 {
		return RaiseErrorVolumeNegative(amount)
	}

	if amount > 100 || as.volumeLevel-amount < 0 {
		amount = as.volumeLevel
	}

	if amount != 0 {
		system.SpeakerLock()
		as.volumeControl.Volume -= float64(amount)

		if as.volumeLevel <= 0 {
			as.volumeControl.Silent = true
		}

		system.SpeakerUnlock()
	}

	return nil
}

// Duration returns the time duration of the
// entire audio stream.
func (as *AudioSource) Duration() time.Duration {
	return as.format.SampleRate.D(as.streamer.Len())
}

// CurrentDuration returns the current position of the
// audio stream in time format.
func (as *AudioSource) CurrentDuration() time.Duration {
	return as.format.SampleRate.D(as.streamer.Position())
}

// Rewind rewinds the audio stream at the specified duration.
func (as *AudioSource) Rewind(t time.Duration) error {
	if t > as.Duration() {
		return RaiseErrorWrongDuration(t)
	}

	pos := as.format.SampleRate.N(t)
	err := as.streamer.Seek(pos)

	if err != nil {
		return err
	}

	return nil
}

// loopIterate plays the audio stream again
// if 'loop' is set to true.
func (as *AudioSource) loopIterate() {
	go func() {
		for {
			select {
			case <-as.loopDone:
				as.loopLocker.Lock()
				loop := as.loop
				as.loopLocker.Unlock()

				if loop {
					system.SpeakerPlay(beep.Seq(as.resampledStreamer, beep.Callback(func() {
						go func() {
							as.loopDone <- true
						}()
					})))
				}

			case <-as.loopCancel:
				return
			}
		}
	}()
}

// Start starts playing the audio stream.
func (as *AudioSource) Start() error {
	system.SpeakerPlay(beep.Seq(as.resampledStreamer, beep.Callback(func() {
		as.loopDone <- true
	})))
	as.loopIterate()

	return nil
}

// Update does nothing.
func (as *AudioSource) Update() error {
	return nil
}

// Destroy pauses the audio stream, closes
// it and then stops the loop goroutine.
func (as *AudioSource) Destroy() error {
	as.Pause()
	err := as.streamer.Close()

	if err != nil {
		return err
	}

	go func() {
		as.loopCancel <- true
	}()

	return nil
}

// AudioSource is a component to play an
// attached sound stream.
func NewAudioSource(name string, audioStream io.ReadCloser) (*AudioSource, error) {
	streamer, format, err := mp3.Decode(audioStream)

	if err != nil {
		return nil, err
	}

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: streamer,
		Base:     1.06,
		Volume:   0,
		Silent:   false,
	}
	resampled := beep.Resample(4,
		format.SampleRate, system.SpeakerSampleRate, streamer)

	as := &AudioSource{
		format:            format,
		streamer:          streamer,
		resampledStreamer: resampled,
		control:           ctrl,
		volumeControl:     volume,
		loop:              false,
		loopLocker:        new(sync.Mutex),
		volumeLevel:       50,
		loopDone:          make(chan bool),
		loopCancel:        make(chan bool),
	}

	as.SetName(name)

	return as, nil
}

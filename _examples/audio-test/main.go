package main

import (
	"os"
	"time"

	"github.com/alacrity-engine/core/audio"
	"github.com/alacrity-engine/core/system"
)

func main() {
	err := system.SpeakerInit()
	handleError(err)

	file, err := os.Open("vyistrel-pistoleta-magnum-357-36128.mp3")
	handleError(err)
	audioSource, err := audio.NewAudioSource("pussy", file)
	handleError(err)

	err = audioSource.Start()
	handleError(err)
	audioSource.SetLoop(true)

	time.Sleep(500 * time.Millisecond)

	newFile, err := os.Open("vyistrel-pistoleta-magnum-357-36128.mp3")
	handleError(err)
	newAudioSource, err := audio.NewAudioSource("pussy", newFile)
	handleError(err)

	err = newAudioSource.Start()
	handleError(err)
	newAudioSource.SetLoop(true)

	time.Sleep(20 * time.Second)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

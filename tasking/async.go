package tasking

import "sync"

const (
	// AsynchronousProcessProgressBufferSize
	// is a capacity of the buffer storing
	// progress notifications about the process.
	AsynchronousProcessProgressBufferSize = 16
)

// AsynchronousProcess is a process
// that takes a lot of time to complete
// and its progress is tracked.
type AsynchronousProcess struct {
	name            string
	progress        int
	locker          *sync.RWMutex
	progressChannel chan int
	result          interface{}
}

// Name returns the name of the
// asynchronous process.
func (ap *AsynchronousProcess) Name() string {
	return ap.name
}

func (ap *AsynchronousProcess) Result() (interface{}, error) {
	ap.locker.RLock()
	defer ap.locker.RUnlock()

	if ap.progress < 100 {
		return nil, NewErrorAsynchronousProcessNotComplete(ap, ap.progress)
	}

	return ap.result, nil
}

func (ap *AsynchronousProcess) SetResult(value interface{}) {
	ap.result = value
}

// ProgressNotifier returns the channel
// to receive notifications about progress
// changes from the asynchronous process.
func (ap *AsynchronousProcess) ProgressNotifier() <-chan int {
	return ap.progressChannel
}

// SetProgress sets the progress
// value for the process.
//
// Should only be called by the
// process initiator.
func (ap *AsynchronousProcess) SetProgress(value int) {
	if value < 0 {
		value = 0
	} else if value > 100 {
		value = 100
	}

	ap.locker.Lock()
	defer ap.locker.Unlock()

	ap.progress = value
	ap.progressChannel <- value
}

// CurrentProgress returns the value
// of the current progress (from 0 to 100)
// of the asynchronous process.
func (ap *AsynchronousProcess) CurrentProgress() int {
	ap.locker.RLock()
	defer ap.locker.RUnlock()

	return ap.progress
}

// NewAsynchronousProcess creates a new
// asynchronous process with the progress
// initially set to 0.
func NewAsynchronousProcess(name string) *AsynchronousProcess {
	return &AsynchronousProcess{
		name:     name,
		progress: 0,
		locker:   new(sync.RWMutex),
		progressChannel: make(chan int,
			AsynchronousProcessProgressBufferSize),
	}
}

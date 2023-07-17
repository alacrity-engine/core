package render

import (
	"runtime"
	"testing"

	"github.com/alacrity-engine/core/system"
	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	length    = 10_000
	locWidth  = 1920
	locHeight = 1080
)

func BenchmarkCopySlice(b *testing.B) {
	slA := make([]byte, length)
	slB := make([]byte, length)

	for i := 0; i < b.N; i++ {
		copy(slA, slB)
	}
}

func BenchmarkCopyBuffer(b *testing.B) {
	runtime.LockOSThread()

	// Initialize the engine.
	b.StopTimer()
	_ = system.InitializeWindow("Demo", locWidth, locHeight, false, false)
	_ = Initialize(locWidth, locHeight, -30, 30)

	var bufA, bufB uint32
	gl.GenBuffers(1, &bufA)
	gl.GenBuffers(1, &bufB)

	slA := make([]byte, length)
	slB := make([]byte, length)

	gl.BindBuffer(gl.ARRAY_BUFFER, bufA)
	gl.BufferData(gl.ARRAY_BUFFER, len(slA), gl.Ptr(slA), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, bufB)
	gl.BufferData(gl.ARRAY_BUFFER, len(slB), gl.Ptr(slB), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.COPY_READ_BUFFER, bufA)
	gl.BindBuffer(gl.COPY_WRITE_BUFFER, bufB)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		gl.CopyBufferSubData(gl.COPY_READ_BUFFER, gl.COPY_WRITE_BUFFER, 0, 0, length)
	}
}

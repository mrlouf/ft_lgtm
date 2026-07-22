package sandbox

import (
	"time"

	"github.com/bytecodealliance/wasmtime-go"
)

// The Sandbox struct stores the WASM runtime engine
// as well as the configuration, ie. the restraints
// on the execution of the WASM code.
type Sandbox struct {
	Runtime *wasmtime.Engine

	MemoryLimit int64

	Timeout time.Duration

	MaxStdout int

	MaxStderr int

	AllowedDirs []string
}

func NewSandbox(memoryLimit int64, timeout time.Duration, maxStdout, maxStderr int, allowedDirs []string) *Sandbox {
	return &Sandbox{
		Runtime:     wasmtime.NewEngine(),
		MemoryLimit: memoryLimit,
		Timeout:     timeout,
		MaxStdout:   maxStdout,
		MaxStderr:   maxStderr,
		AllowedDirs: allowedDirs,
	}
}

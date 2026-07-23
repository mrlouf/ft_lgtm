package sandbox

import (
	"context"
	"lgtm/internal/test"
	"log"
	"time"

	"github.com/tetratelabs/wazero"
)

type WazeroSandbox struct {
	runtime     wazero.Runtime
	memoryLimit int64
	timeout     time.Duration
	maxStdout   int
	maxStderr   int
	allowedDirs []string
}

type Option func(*WazeroSandbox)

func WithTimeout(d time.Duration) Option {
	return func(s *WazeroSandbox) { s.timeout = d }
}

func WithMemoryLimit(bytes int64) Option {
	return func(s *WazeroSandbox) { s.memoryLimit = bytes }
}

func NewWazeroSandbox(opts ...Option) *WazeroSandbox {

	s := &WazeroSandbox{
		memoryLimit: 64 * 1024 * 1024,
		timeout:     10 * time.Second,
		maxStdout:   1024 * 1024,
		maxStderr:   1024 * 1024,
		allowedDirs: []string{"/tmp"},
	}

	// set options if provided to override defaults
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *WazeroSandbox) Compile(ctx context.Context, source []byte) ([]byte, error) {

	log.Println("compile: start")

	if err := test.SleepOrCancel(ctx, 2*time.Second, "compile"); err != nil {
		return nil, err
	}

	log.Println("compile: done")

	return []byte("fake-wasm-binary"), nil

}

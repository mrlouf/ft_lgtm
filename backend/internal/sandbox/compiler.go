package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func compileGoToWasm(ctx context.Context, source []byte) ([]byte, error) {

	tmpDir, err := os.MkdirTemp("", "snippet-*")
	if err != nil {
		return nil, fmt.Errorf("compile: tmpdir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(srcPath, source, 0o644); err != nil {
		return nil, fmt.Errorf("compile: write source: %w", err)
	}

	outPath := filepath.Join(tmpDir, "out.wasm")
	cmd := exec.CommandContext(ctx, "tinygo", "build", "-o", outPath, "-target=wasi", srcPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("compile: %s: %w", stderr.String(), err)
	}

	// * DEBUG
	log.Printf("compile: done, wasm binary at %s", outPath)
	log.Printf("compile: binary size: %d bytes", func() int64 {
		info, err := os.Stat(outPath)
		if err != nil {
			return 0
		}
		return info.Size()
	}())

	return os.ReadFile(outPath)
}

func (s *WazeroSandbox) Compile(ctx context.Context, source []byte, lang string) ([]byte, error) {

	log.Println("compile: start")

	var wasmBinary []byte
	var err error

	switch lang {
	case "go":
		wasmBinary, err = compileGoToWasm(ctx, source)
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	log.Println("compile: done")

	return wasmBinary, err

}

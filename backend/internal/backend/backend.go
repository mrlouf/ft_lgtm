package backend

import (
	"context"
	"fmt"
)

type Compiler interface {
	Compile(ctx context.Context, source []byte) ([]byte, error)
}

type Executor interface {
	Execute(ctx context.Context, wasmBinary []byte) (stdout, stderr string, err error)
}

type Publisher interface {
	Publish(ctx context.Context, data []byte) (cid string, err error)
}

type Backend struct {
	Compiler  Compiler
	Executor  Executor
	Publisher Publisher
}

func (b *Backend) Run(ctx context.Context, source []byte) (string, string, error) {

	wasmBinary, err := b.Compiler.Compile(ctx, source)
	if err != nil {
		return "", "", fmt.Errorf("compile: %w", err)
	}

	return b.Executor.Execute(ctx, wasmBinary)
}

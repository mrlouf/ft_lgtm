package backend

import (
	"context"
	"fmt"
	"log"
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

func NewBackend(compiler Compiler, executor Executor, publisher Publisher) *Backend {
	return &Backend{
		Compiler:  compiler,
		Executor:  executor,
		Publisher: publisher,
	}
}

func (b *Backend) Run(ctx context.Context, source []byte, language string) (string, string, string, error) {

	log.Printf("Run: start for language: %s", language)

	wasmBinary, err := b.Compiler.Compile(ctx, source)
	if err != nil {
		return "", "", "", fmt.Errorf("compile: %w", err)
	}

	stdout, stderr, err := b.Executor.Execute(ctx, wasmBinary)
	if err != nil {
		return "", "", "", fmt.Errorf("execute: %w", err)
	}

	cid, err := b.Publisher.Publish(ctx, []byte(stdout))
	if err != nil {
		return "", "", "", fmt.Errorf("publish: %w", err)
	}

	return stdout, stderr, cid, nil

}

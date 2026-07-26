package backend

import (
	"context"
	"lgtm/internal/ipfs"
	"log"
)

type Compiler interface {
	Compile(ctx context.Context, source []byte, lang string) ([]byte, error)
}

type Executor interface {
	Execute(ctx context.Context, wasmBinary []byte) (stdout, stderr string, err error)
}

type Publisher interface {
	Publish(ctx context.Context, source []byte, stdout []byte) (response ipfs.ResponseCID, err error)
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

func (b *Backend) Run(ctx context.Context, source []byte, language string) (string, string, ipfs.ResponseCID, error) {

	log.Printf("Run: start for language: %s", language)

	wasmBinary, err := b.Compiler.Compile(ctx, source, language)
	if err != nil {
		return "", "", ipfs.ResponseCID{}, err
	}

	stdout, stderr, err := b.Executor.Execute(ctx, wasmBinary)
	if err != nil {
		return stdout, stderr, ipfs.ResponseCID{}, err
	}

	responseCID, err := b.Publisher.Publish(ctx, source, []byte(stdout))
	if err != nil {
		return stdout, stderr, ipfs.ResponseCID{}, err
	}

	return stdout, stderr, responseCID, nil

}

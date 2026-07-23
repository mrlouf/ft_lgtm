package sandbox

import (
	"context"
	"log"
)

type WazeroExecutor struct {
	*WazeroSandbox
}

func NewWazeroExecutor(sandbox *WazeroSandbox) *WazeroExecutor {
	return &WazeroExecutor{
		WazeroSandbox: sandbox,
	}
}

func (e *WazeroExecutor) Execute(ctx context.Context, wasmBinary []byte) (stdout, stderr string, err error) {

	log.Println("execute: start")

	log.Println("execute: ", string(wasmBinary))

	log.Println("execute: done")

	return "", "", nil

}

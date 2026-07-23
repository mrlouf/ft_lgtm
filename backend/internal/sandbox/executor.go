package sandbox

import (
	"context"
	"lgtm/internal/test"
	"log"
	"time"
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

	if err := test.SleepOrCancel(ctx, 2*time.Second, "execute"); err != nil {
		return "", "", err
	}

	log.Println("execute: done")

	return "", "", nil

}

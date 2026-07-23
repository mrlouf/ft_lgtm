package sandbox

import "context"

type WazeroExecutor struct {
	*WazeroSandbox
}

func NewWazeroExecutor(sandbox *WazeroSandbox) *WazeroExecutor {
	return &WazeroExecutor{
		WazeroSandbox: sandbox,
	}
}

func (e *WazeroExecutor) Execute(ctx context.Context, wasmBinary []byte) (stdout, stderr string, err error) {

	return "", "", nil

}

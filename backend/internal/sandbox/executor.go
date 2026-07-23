package sandbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

type WazeroExecutor struct {
	runtime wazero.Runtime
}

func NewWazeroExecutor(ctx context.Context) *WazeroExecutor {
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	return &WazeroExecutor{runtime: r}
}

func (e *WazeroExecutor) Execute(ctx context.Context, wasmBytes []byte) (stdout, stderr string, err error) {

	log.Println("execute: start")

	var stdoutBuf, stderrBuf bytes.Buffer

	config := wazero.NewModuleConfig().
		WithStdout(&stdoutBuf).
		WithStderr(&stderrBuf)

	// Compile the WebAssembly module from the raw bytes
	compiled, err := e.runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		return "", "", fmt.Errorf("compile module: %w", err)
	}
	defer compiled.Close(ctx)

	// Instantiate the module with the configuration
	// aka run the WebAssembly module
	mod, err := e.runtime.InstantiateModule(ctx, compiled, config)
	if mod != nil {
		defer mod.Close(ctx)
	}
	if err != nil {
		var exitErr *sys.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 0 {
			// proc_exit(0): The program exited successfully,
			// but we still want to capture stdout and stderr.
			log.Println("execute: done (clean exit)")
			return stdoutBuf.String(), stderrBuf.String(), nil
		}
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("instantiate: %w", err)
	}

	log.Println("execute: done")
	return stdoutBuf.String(), stderrBuf.String(), nil

}

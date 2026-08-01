package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"go.opentelemetry.io/otel/trace"
)

type WazeroExecutor struct {
	runtime wazero.Runtime
	tracer  trace.Tracer
	logger  *slog.Logger
}

func NewWazeroExecutor(ctx context.Context, t trace.Tracer, l *slog.Logger) *WazeroExecutor {
	r := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return &WazeroExecutor{
		runtime: r,
		tracer:  t,
		logger:  l,
	}
}

func (e *WazeroExecutor) Execute(ctx context.Context, wasmBytes []byte) (stdout, stderr string, err error) {

	e.logger.InfoContext(ctx, "execute: start")
	ctx, span := e.tracer.Start(ctx, "executor.execute")
	defer span.End()

	var stdoutBuf, stderrBuf bytes.Buffer

	config := wazero.NewModuleConfig().
		WithStdout(&stdoutBuf).
		WithStderr(&stderrBuf)

	// Compile the WebAssembly module from the raw bytes
	compiled, err := e.runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		e.logger.ErrorContext(ctx, "execute: error compiling module", "error", err)
		span.RecordError(err)
		return "", "", fmt.Errorf("compile module: %w", err)
	}
	defer compiled.Close(ctx)

	e.logger.InfoContext(ctx, "execute: module compiled")
	span.AddEvent("module compiled")

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
			e.logger.InfoContext(ctx, "execute: done (clean exit)")
			span.AddEvent("module executed (clean exit)")
			return stdoutBuf.String(), stderrBuf.String(), nil
		}
		e.logger.ErrorContext(ctx, "execute: error instantiating module", "error", err)
		span.RecordError(err)
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("instantiate: %w", err)
	}

	// proc_exit(0): The program exited successfully,
	// but we still want to capture stdout and stderr.
	e.logger.InfoContext(ctx, "execute: done (clean exit)")
	span.AddEvent("module executed (clean exit)")
	return stdoutBuf.String(), stderrBuf.String(), nil
}

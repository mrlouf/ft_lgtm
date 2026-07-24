package backend

import (
	"context"
	"errors"
	"fmt"
)

type CompileError struct{ Err error }

func (e *CompileError) Error() string { return fmt.Sprintf("compile: %v", e.Err) }
func (e *CompileError) Unwrap() error { return e.Err }

type ExecuteError struct{ Err error }

func (e *ExecuteError) Error() string { return fmt.Sprintf("execute: %v", e.Err) }
func (e *ExecuteError) Unwrap() error { return e.Err }

type PublishError struct{ Err error }

func (e *PublishError) Error() string { return fmt.Sprintf("publish: %v", e.Err) }
func (e *PublishError) Unwrap() error { return e.Err }

type RunStatus int

const (
	StatusUnknown RunStatus = iota
	StatusSuccess
	StatusCompileError
	StatusExecuteError
	StatusTimeout
)

func (s RunStatus) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusCompileError:
		return "compile_error"
	case StatusExecuteError:
		return "execute_error"
	case StatusTimeout:
		return "timeout"
	default:
		return "unknown"
	}
}

func ClassifyError(err error) RunStatus {
	var compileErr *CompileError
	var executeErr *ExecuteError
	switch {
	case err == nil:
		return StatusSuccess
	case errors.Is(err, context.DeadlineExceeded):
		return StatusTimeout
	case errors.As(err, &compileErr):
		return StatusCompileError
	case errors.As(err, &executeErr):
		return StatusExecuteError
	default:
		return StatusUnknown
	}
}

package test

import (
	"context"
	"fmt"
	"time"
)

// sleepOrCancel sleeps for the specified duration or returns early if the context is canceled.
// This is useful to simulate long-running operations in the different stages
// of the backend pipeline (compile, execute, publish) and test cancellation.
func SleepOrCancel(ctx context.Context, d time.Duration, stage string) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: interrupted: %w", stage, ctx.Err())
	}
}

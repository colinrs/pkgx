package contextx

import (
	"context"
	"time"
)

// ShrinkDeadline returns a new Context with proper deadline base on the given ctx and timeout.
// And returns a cancel function as well.
func ShrinkDeadline(ctx context.Context, timeout time.Duration) (context.Context, func()) {
	if deadline, ok := ctx.Deadline(); ok {
		leftTime := time.Until(deadline)
		if leftTime < timeout {
			timeout = leftTime
		}
	}

	return context.WithDeadline(ctx, time.Now().Add(timeout))
}

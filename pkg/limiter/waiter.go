package limiter

import (
	"context"

	"golang.org/x/time/rate"
)

type waiter interface {
	UnLimit() bool
	WaitN(context.Context, int) error
	Limit() rate.Limit
	Take(n int)
	SetLimit(newLimit int32)
}

type BaseWaiter struct {
	*rate.Limiter
}

func NewBaseWaiter(b int) *BaseWaiter {
	currentLimit := rate.Limit(b) * 1024
	if currentLimit <= 0 {
		currentLimit = rate.Inf
	}

	return &BaseWaiter{
		Limiter: rate.NewLimiter(currentLimit, limiterBurstSize),
	}
}

func (w *BaseWaiter) UnLimit() bool {
	return w.Limit() == rate.Inf
}

func (w *BaseWaiter) WaitN(ctx context.Context, n int) error {
	return w.Limiter.WaitN(ctx, n)
}

func (w *BaseWaiter) Limit() rate.Limit {
	return w.Limiter.Limit()
}

func (w *BaseWaiter) SetLimit(newLimit int32) {
	currentLimit := rate.Limit(newLimit) * 1024
	if currentLimit <= 0 {
		currentLimit = rate.Inf
	}

	w.Limiter.SetLimit(currentLimit)
}

func (w *BaseWaiter) Take(n int) {
	if n < limiterBurstSize {
		_ = w.WaitN(context.TODO(), n)
		return
	}

	for n > 0 {
		if n > limiterBurstSize {
			_ = w.WaitN(context.TODO(), n)
			n -= limiterBurstSize
		} else {
			_ = w.WaitN(context.TODO(), n)
			n = 0
		}
	}
}

package limiter

import (
	"context"
	
	"golang.org/x/time/rate"
)

type waiter interface {
	WaitN(context.Context, int) error
	Limit() rate.Limit
}

type baseWaiter rate.Limiter

func (w *baseWaiter) UnLimit() bool {
	return (*rate.Limiter)(w).Limit() == rate.Inf
}

func (w *baseWaiter) Limit() rate.Limit {
	return (*rate.Limiter)(w).Limit()
}

func (w *baseWaiter) take(n int) {
	l := (*rate.Limiter)(w)
	if n < limiterBurstSize {
		_ = l.WaitN(context.TODO(), n)
		return
	}
	
	for n > 0 {
		if n > limiterBurstSize {
			_ = l.WaitN(context.TODO(), n)
			n -= limiterBurstSize
		} else {
			_ = l.WaitN(context.TODO(), n)
			n = 0
		}
	}
}

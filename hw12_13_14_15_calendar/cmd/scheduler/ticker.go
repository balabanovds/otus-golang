package main

import (
	"context"
	"time"
)

type worker func(context.Context, time.Time)

func every(ctx context.Context, d time.Duration, work worker) {
	ticker := time.NewTicker(d)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				work(ctx, t)
			case <-ctx.Done():
				return
			}
		}
	}()
}

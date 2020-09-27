package main

import (
	"context"
	"time"
)

type worker func(context.Context, time.Time)

func every(ctx context.Context, done chan struct{}, interval time.Duration, worker worker) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			done <- struct{}{}
			return
		case t := <-ticker.C:
			worker(ctx, t)
		}
	}
}

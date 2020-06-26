package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type counter struct {
	sync.Mutex
	threshold    int
	current      int
	ignoreErrors bool
}

func newCounter(threshold int) *counter {
	return &counter{
		threshold:    threshold,
		ignoreErrors: threshold <= 0,
	}
}

func (c *counter) increment() {
	c.Lock()
	defer c.Unlock()
	c.current++
}

func (c *counter) reachedThreshold() bool {
	c.Lock()
	defer c.Unlock()
	return !c.ignoreErrors && c.current >= c.threshold
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	ctr := newCounter(m)

	var wg sync.WaitGroup
	wg.Add(n)

	tasksCh := make(chan Task)

	// we create pool of n workers
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			// consuming chan for Task unless all done or chan closed by main goroutine
			for task := range tasksCh {
				if task() != nil {
					// if task return error we increment counter
					ctr.increment()
				}
			}
		}()
	}

	// iterating over tasks and producing then to chan for workers pool
	for _, task := range tasks {
		if ctr.reachedThreshold() {
			break
		}
		tasksCh <- task
	}

	close(tasksCh)

	wg.Wait()

	if ctr.reachedThreshold() {
		return ErrErrorsLimitExceeded
	}
	return nil
}

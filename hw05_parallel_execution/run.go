package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var nerrors atomic.Int32
	var wg sync.WaitGroup

	queue := make(chan Task, len(tasks))

	for _, t := range tasks {
		queue <- t
	}

	close(queue)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() error {
			defer wg.Done()

			for task := range queue {
				if err := task(); err != nil {
					nerrors.Add(1)
				}

				if m > 0 && int(nerrors.Load()) >= m {
					return ErrErrorsLimitExceeded
				}
			}
			return nil
		}()
	}

	wg.Wait()

	if m > 0 && int(nerrors.Load()) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

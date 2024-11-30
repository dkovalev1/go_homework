package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type taskQueue struct {
	mtx   sync.Mutex
	tasks []Task
}

func (q *taskQueue) getNewTask() Task {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.tasks) == 0 {
		return nil
	}

	task := q.tasks[0]
	q.tasks = q.tasks[1:]

	return task
}

func newTaskQueue(tasks []Task) *taskQueue {
	return &taskQueue{
		tasks: tasks,
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var nerrors atomic.Int32
	var wg sync.WaitGroup

	queue := newTaskQueue(tasks)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() error {
			defer wg.Done()

			for task := queue.getNewTask(); task != nil; task = queue.getNewTask() {
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

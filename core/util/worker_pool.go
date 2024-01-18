package util

import (
	"sync"
)

type WorkerPool struct {
	workers int
	wg      sync.WaitGroup
	errc    chan error
	errors  []error
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{workers: workers, errc: make(chan error), errors: make([]error, workers)}
}

func (pool *WorkerPool) Go(fn func(worker int) error) {
	pool.wg.Add(pool.workers)
	for i := 0; i < pool.workers; i++ {
		go func(worker int) {
			if err := fn(worker); err != nil {
				pool.errc <- err
			}
			pool.wg.Done()
		}(i)
	}
	go func() {
		for err := range pool.errc {
			pool.errors = append(pool.errors, err)
		}
	}()
}

func (pool *WorkerPool) Wait(fn func()) {
	go func() {
		pool.wg.Wait()
		close(pool.errc)
		fn()
	}()
}

func (pool *WorkerPool) HasErrors() bool {
	return len(pool.errors) > 0
}

func (pool *WorkerPool) Errors() []error {
	return pool.errors
}

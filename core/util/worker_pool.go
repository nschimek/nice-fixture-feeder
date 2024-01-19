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
	return &WorkerPool{workers: workers, errc: make(chan error), errors: make([]error, 0, workers)}
}

// Go runs the given function on one of the available workers, if any.
// Non-nil Errors are collected and available after execution via HasErrors() and Errors().
func (pool *WorkerPool) Go(fn func(worker int) error) {
	pool.wg.Add(pool.workers)
	for i := 0; i < pool.workers; i++ {
		i := i
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

// Wait blocks until the worker pool completes its tasks, then runs the given function, if any.
func (pool *WorkerPool) Wait(fn func()) {
	go func() {
		pool.wg.Wait()
		close(pool.errc)
		if fn != nil {
			fn()
		}
	}()
}

// HasErrors returns whether or not erorrs were output by the Go function.
func (pool *WorkerPool) HasErrors() bool {
	return len(pool.errors) > 0
}

// Errors returns the collected errors.
func (pool *WorkerPool) Errors() []error {
	return pool.errors
}

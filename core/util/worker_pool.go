package util

import (
	"sync"
)

type WorkerPool struct {
	size int
	wg sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{size: size}
}

func (pool *WorkerPool) Go(fn func(worker int)) {
	pool.wg.Add(pool.size)
	for i := 0; i < pool.size; i++ {
		go func(worker int) {
			fn(worker)
			pool.wg.Done()
		}(i)
	}
}

func (pool *WorkerPool) Wait(fn func()) {
	go func() {
		pool.wg.Wait()
		fn()
	}()
}
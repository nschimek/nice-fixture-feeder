package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	w := 0
	in := make(chan int)
	out := make(chan int)

	pool := NewWorkerPool(3)
	pool.Go(func(worker int) error {
		for range in {
			time.Sleep(time.Duration(100 * time.Millisecond))
			out <- worker
		}
		return nil
	})
	pool.Wait(func() {
		w++
		close(out)
	})

	go func() {
		defer close(in)
		for i := 0; i < 9; i++ {
			in <- i
		}
	}()

	m := make(map[int]int) // keep track of runs per worker
	for o := range out {
		m[o]++
	}

	assert.Len(t, m, 3)
	assert.Equal(t, 1, w)
	assert.Equal(t, m[0], 3)
	assert.Equal(t, m[1], 3)
	assert.Equal(t, m[2], 3)
	assert.False(t, pool.HasErrors())
}

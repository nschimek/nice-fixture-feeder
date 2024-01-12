package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdMapToArray(t *testing.T) {
	m := map[int]struct{}{
		1: Exists,
		2: Exists,
	}
	e := []int{1, 2}
	a := IdMapToArray(m)

	assert.ElementsMatch(t, e, a)
}

func TestIdArrayToMap(t *testing.T) {
	e := []int{1, 2}
	a := IdArrayToMap(e)

	assert.Equal(t, 2, len(e))
	assert.Contains(t, a, e[0])
	assert.Contains(t, a, e[1])
}

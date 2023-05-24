package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdMapToArray(t *testing.T) {
	m := map[string]struct{}{
		"test1": exists,
		"test2": exists,
	}
	e := IdMapToArray(m)
	assert.Equal(t, 2, len(e))

	_, ok1 := m[e[0]]
	assert.True(t, ok1)

	_, ok2 := m[e[1]]
	assert.True(t, ok2)
}

func TestIdArrayToMap(t *testing.T) {
	a := []string{"test1", "test2"}
	e := IdArrayToMap(a)

	assert.Equal(t, 2, len(e))

	_, ok1 := e[a[0]]
	assert.True(t, ok1)

	_, ok2 := e[a[1]]
	assert.True(t, ok2)
}
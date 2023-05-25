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
	e := []string{"test1", "test2"}
	a := IdMapToArray(m)
	
	assert.ElementsMatch(t, e, a)
}

func TestIdArrayToMap(t *testing.T) {
	e := []string{"test1", "test2"}
	a := IdArrayToMap(e)

	assert.Equal(t, 2, len(e))
	assert.Contains(t, a, e[0])
	assert.Contains(t, a, e[1])
}
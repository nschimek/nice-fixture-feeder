package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapOrEmptyArrayUnmarshalValid(t *testing.T) {
	str := `{"field 1": "test 2", "field 2": "test 2"}`
	m := map[string]string{"field 1": "test 2", "field 2": "test 2"}

	e := new(MapOrEmptyArray)
	*e = MapOrEmptyArray(m)

	a := new(MapOrEmptyArray)
	err := a.UnmarshalJSON([]byte(str))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestMapOrEmptyArrayUnmarshalEmpty(t *testing.T) {
	e := new(MapOrEmptyArray)
	*e = MapOrEmptyArray(nil)

	a := new(MapOrEmptyArray)
	err := a.UnmarshalJSON([]byte("[]"))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestMapOrEmptyArrayUnmarshalError(t *testing.T) {
	a := new(MapOrEmptyArray)
	err := a.UnmarshalJSON([]byte(`{invalidJson}`))

	assert.ErrorContains(t, err, "invalid character")
}
package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinnerBoolValidW(t *testing.T) {
	v := "true"

	e := new(WinnerBool)
	*e = WinnerBool("W")

	a := new(WinnerBool)
	err := a.UnmarshalJSON([]byte(v))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestWinnerBoolValidL(t *testing.T) {
	v := "false"

	e := new(WinnerBool)
	*e = WinnerBool("L")

	a := new(WinnerBool)
	err := a.UnmarshalJSON([]byte(v))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestWinnerBoolValidT(t *testing.T) {
	v := "null"

	e := new(WinnerBool)
	*e = WinnerBool("D")

	a := new(WinnerBool)
	err := a.UnmarshalJSON([]byte(v))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestWinnerBoolInvalid(t *testing.T) {
	a := new(WinnerBool)
	assert.Error(t, a.UnmarshalJSON([]byte("ASDF")))
}
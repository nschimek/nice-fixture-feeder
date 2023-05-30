package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalValid(t *testing.T) {
	date := "2023-05-26"
	dateTime, _ := time.Parse("2006-01-02", date)
	
	e := new(CivilTime)
	*e = CivilTime(dateTime)

	a := new(CivilTime)
	a.UnmarshalJSON([]byte(date))

	assert.Equal(t, e, a)
}

func TestUnmarshalNull(t *testing.T) {
	e := new(CivilTime)
	*e = CivilTime(time.Time{})

	a := new(CivilTime)
	a.UnmarshalJSON([]byte(""))

	assert.Equal(t, e, a)
}

func TestUnmarshalInvalid(t *testing.T) {
	a := new(CivilTime)
	assert.Error(t, a.UnmarshalJSON([]byte("20220529")))
}

func TestMarshal(t *testing.T) {
	date := "2023-05-26"
	dateTime, _ := time.Parse("2006-01-02", date)

	a := new(CivilTime)
	*a = CivilTime(dateTime)
	ab, _ := a.MarshalJSON() // actual bytes

	eb := []byte(fmt.Sprintf("\"%s\"", date)) // expected bytes
	
	assert.Equal(t, eb, ab)
}
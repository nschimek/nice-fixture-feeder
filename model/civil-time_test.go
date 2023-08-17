package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCivilTimeUnmarshalValid(t *testing.T) {
	date := "2023-05-26"
	dateTime, _ := time.Parse("2006-01-02", date)
	
	e := new(CivilTime)
	*e = CivilTime(dateTime)

	a := new(CivilTime)
	err := a.UnmarshalJSON([]byte(date))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestCivilTimeUnmarshalNull(t *testing.T) {
	e := new(CivilTime)
	*e = CivilTime(time.Time{})

	a := new(CivilTime)
	err := a.UnmarshalJSON([]byte(""))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestCivilTimeUnmarshalInvalid(t *testing.T) {
	a := new(CivilTime)
	assert.Error(t, a.UnmarshalJSON([]byte("20220529")))
}

func TestCivilTimeMarshal(t *testing.T) {
	date := "2023-05-26"
	dateTime, _ := time.Parse("2006-01-02", date)

	a := new(CivilTime)
	*a = CivilTime(dateTime)
	ab, _ := a.MarshalJSON() // actual bytes

	eb := []byte(fmt.Sprintf("\"%s\"", date)) // expected bytes
	
	assert.Equal(t, eb, ab)
}

func TestCivilTimeScanValid(t *testing.T) {
	date := "2023-05-26"
	dateTime, _ := time.Parse("2006-01-02", date)
	
	e := new(CivilTime)
	*e = CivilTime(dateTime)

	a := new(CivilTime)
	err := a.Scan(dateTime)

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestCivilTimecanInvalid(t *testing.T) {
	a := new(CivilTime)
	assert.Error(t, a.Scan("asdf"))
}

func TestCivilTimeValue(t *testing.T) {
	date := "2023-05-26"
	dateTime, _ := time.Parse("2006-01-02", date)

	a := new(CivilTime)
	*a = CivilTime(dateTime)
	ad, _ := a.Value()

	assert.Equal(t, date, ad)
}
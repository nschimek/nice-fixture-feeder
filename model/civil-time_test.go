package model

import (
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
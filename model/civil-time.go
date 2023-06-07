package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
)

type CivilTime time.Time

// JSON Unmarshal and Marshal interface implementations
func (c *CivilTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
			return nil
	}

	t, err := time.Parse(core.YYYY_MM_DD, value) //parse time
	if err != nil {
			return err
	}
	
	*c = CivilTime(t) //set result using the pointer
	return nil
}

func (c CivilTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format(core.YYYY_MM_DD) + `"`), nil
}

// Scanner / Value interface methods to integrate with GORM
func (c *CivilTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to convert DateTime value:", value))
	}
	*c = CivilTime(t)
	return nil
}

func (c CivilTime) Value() (driver.Value, error) {
	return time.Time(c).Format(core.YYYY_MM_DD), nil
} 

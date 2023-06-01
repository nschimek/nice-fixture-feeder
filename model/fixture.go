package model

import (
	"errors"
	"time"
)

type Fixture struct {

}

type FixtureFixture struct {
	Id int `gorm:"primaryKey"`
	Date time.Time
	Venue FixtureVenue `gorm:"embedded;embeddedPrefix:venue_"`
	Status FixtureStatus
	League FixtureLeague `gorm:"embedded;embeddedPrefix:league_"`
	Teams FixtureTeams `gorm:"embedded;embeddedPrefix:team_"`
	Goals FixtureGoals `gorm:"embedded;embeddedPrefix:goals_"`
}

type FixtureLeague struct {
	Id int
}

type FixtureVenue struct {
	Name, City string
}

type FixtureStatus struct {
	Id string `json:"short" gorm:"primaryKey"`
	Name string `json:"long"`
	Type string `json:"-"`
	Description string `json:"-"`
}

type FixtureTeams struct {
	Home FixtureTeam `gorm:"embedded;embeddedPrefix:home_"`
	Away FixtureTeam `gorm:"embedded;embeddedPrefix:away_"`
}

type FixtureTeam struct {
	Id int
	Winner WinnerBool
}

type FixtureGoals struct {
	Home, Away int
}

type WinnerBool string

// converts the API's true/false/null winner field values to W/L/D
func (w *WinnerBool) UnmarshalJSON(data []byte) error {
	v := string(data)
	var s string

	if v == "null" {
		s = "D"
	} else if v == "true" {
		s = "W"
	} else if v == "false" {
		s = "L"
	} else {
		return errors.New("unexpected value in winner boolean field")
	}

	*w = WinnerBool(s)
	return nil
}
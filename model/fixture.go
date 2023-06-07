package model

import (
	"errors"
	"time"
)

type Fixture struct {
	Fixture FixtureFixture `gorm:"embedded"`
	League FixtureLeague `gorm:"embedded"`
	Teams FixtureTeams `gorm:"embedded;embeddedPrefix:team_"`
	Goals FixtureGoals `gorm:"embedded;embeddedPrefix:goals_"`
}

func (f *Fixture) GetTeamStatsId(home bool) TeamStatsId {
	// sometimes the lack of a ternary operator is annoying...
	t := f.Teams.Home.Id
	if !home {
		t = f.Teams.Away.Id
	}

	return TeamStatsId{
		TeamId: t,
		LeagueId: f.League.Id,
		Season: f.League.Season,
		FixtureId: f.Fixture.Id,
		Date: f.Fixture.Date.AddDate(0, 0, -1),
	}
}

type FixtureFixture struct {
	Id int `gorm:"primaryKey"`
	Date time.Time
	Venue FixtureVenue `gorm:"embedded;embeddedPrefix:venue_"`
	Status FixtureStatus `gorm:"embedded;embeddedPrefix:status_"`
}

type FixtureLeague struct {
	Id int `gorm:"column:league_id"`
	Season int
}

type FixtureVenue struct {
	Name, City string
}

type FixtureStatus struct {
	Id string `json:"short"`
}

type FixtureTeams struct {
	Home FixtureTeam `gorm:"embedded;embeddedPrefix:home_"`
	Away FixtureTeam `gorm:"embedded;embeddedPrefix:away_"`
}

type FixtureTeam struct {
	Id int
	Result WinnerBool `json:"winner"`
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
package model

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

const (
	ResultWin = "W"
	ResultLoss = "L"
	ResultDraw = "D"
)

type Fixture struct {
	Fixture FixtureFixture `gorm:"embedded"`
	League FixtureLeague `gorm:"embedded"`
	Teams FixtureTeams `gorm:"embedded;embeddedPrefix:team_"`
	Goals FixtureGoals `gorm:"embedded;embeddedPrefix:goals_"`
	Audit   `json:"-"`
}

// not persisted, but used for maintaining stats
type FixtureResultStats struct {
	Home bool
	Result string
	GoalsFor, GoalsAgainst int
}

func (f *Fixture) GetTeamStatsId(home bool) *TeamStatsId {
	// sometimes the lack of a ternary operator is annoying...
	t := f.Teams.Home.Id
	if !home {
		t = f.Teams.Away.Id
	}

	return &TeamStatsId{
		TeamId: t,
		LeagueId: f.League.Id,
		Season: f.League.Season,
		FixtureId: f.Fixture.Id,
	}
}

// Get an instance of TSID from the Fixture but populate NextFixtureID with the ID instead
func (f *Fixture) GetTeamStatsNextId(home bool) *TeamStatsId {
	tsid := f.GetTeamStatsId(home)
	return &TeamStatsId{
		TeamId: tsid.TeamId,
		LeagueId: tsid.LeagueId,
		Season: tsid.Season,
		NextFixtureId: tsid.FixtureId,
	}
}

func (f *Fixture) GetResultStats(teamId int) FixtureResultStats {
	if (f.Teams.Home.Id == teamId) {
		return FixtureResultStats{
			Home: true,
			Result: string(f.Teams.Home.Result),
			GoalsFor: f.Goals.Home,
			GoalsAgainst: f.Goals.Away,
		}
	} else {
		return FixtureResultStats{
			Home: false,
			Result: string(f.Teams.Away.Result),
			GoalsFor: f.Goals.Away,
			GoalsAgainst: f.Goals.Home,
		}
	}
}


type FixtureFixture struct {
	Id int `gorm:"primaryKey"`
	Date time.Time
	Venue FixtureVenue `gorm:"embedded;embeddedPrefix:venue_"`
	Status FixtureStatusId `gorm:"embedded;embeddedPrefix:status_"`
}

type FixtureLeague struct {
	Id int `gorm:"column:league_id"`
	Season int
	Round Round 
}

type FixtureVenue struct {
	Name, City string
}

type FixtureStatusId struct {
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
		s = ResultDraw
	} else if v == "true" {
		s = ResultWin
	} else if v == "false" {
		s = ResultLoss
	} else {
		return errors.New("unexpected value in winner boolean field")
	}

	*w = WinnerBool(s)
	return nil
}

type Round int

// converts the API's round string value into an integer
func (r *Round) UnmarshalJSON(data []byte) error {
	v := string(data)
	var i int

	if v == "null" {
		i = 0
	} else {
		re := regexp.MustCompile(`[^0-9]+`)
		i, _ = strconv.Atoi(re.ReplaceAllString(v, ""))
	}

	*r = Round(i)
	return nil
}
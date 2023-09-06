package model

import (
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/stretchr/testify/assert"
)

func TestGetTeamStatsId(t *testing.T) {
	f := Fixture{
		Fixture: FixtureFixture{
			Id: 100,
			Date: time.Date(2023, 3, 5, 16, 30, 0, 0, core.UTC),
		},
		League: FixtureLeague{Id: 39, Season: 2022},
		Teams: FixtureTeams{
			Home: FixtureTeam{Id: 40},
			Away: FixtureTeam{Id: 33},
		},
		Goals: FixtureGoals{Home: 7, Away: 0},
	}
	h := f.GetTeamStatsId(true)
	a := f.GetTeamStatsId(false)
	p := f.GetTeamStatsIdPrevSeason(true)

	assert.Equal(t, f.Teams.Home.Id, h.TeamId)
	assert.Equal(t, f.Teams.Away.Id, a.TeamId)
	assert.Equal(t, f.Fixture.Id, h.FixtureId)
	assert.Equal(t, f.Fixture.Id, a.FixtureId)
	assert.Equal(t, f.League.Id, h.LeagueId)
	assert.Equal(t, f.League.Id, a.LeagueId)
	assert.Equal(t, f.League.Season, h.Season)
	assert.Equal(t, f.League.Season, a.Season)
	assert.Equal(t, f.League.Season - 1, p.Season)
}

func TestGetResultStats(t *testing.T) {
	f := Fixture{
		Fixture: FixtureFixture{
			Id: 100,
			Date: time.Date(2023, 3, 5, 16, 30, 0, 0, core.UTC),
		},
		League: FixtureLeague{Id: 39, Season: 2022},
		Teams: FixtureTeams{
			Home: FixtureTeam{Id: 40, Result: "W"},
			Away: FixtureTeam{Id: 33, Result: "L"},
		},
		Goals: FixtureGoals{Home: 7, Away: 1},
	}
	lp := f.GetResultStats(40)
	mu := f.GetResultStats(33)

	assert.Equal(t, lp.GoalsFor, 7)
	assert.Equal(t, lp.GoalsAgainst, 1)
	assert.True(t, lp.Home)
	assert.Equal(t, "W", lp.Result)
	assert.Equal(t, mu.GoalsFor, 1)
	assert.Equal(t, mu.GoalsAgainst, 7)
	assert.False(t, mu.Home)
	assert.Equal(t, "L", mu.Result)
}

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

func TestResult(t *testing.T) {
	v := "Regular Season - 5"

	e := new(Round)
	*e = Round(5)

	a := new(Round)
	err := a.UnmarshalJSON([]byte(v))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestResultNull(t *testing.T) {
	v := "null"

	e := new(Round)
	*e = Round(0)

	a := new(Round)
	err := a.UnmarshalJSON([]byte(v))

	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

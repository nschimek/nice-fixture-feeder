package model

import (
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
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

	assert.Equal(t, f.Teams.Home.Id, h.TeamId)
	assert.Equal(t, f.Teams.Away.Id, a.TeamId)
	assert.Equal(t, f.Fixture.Id, h.FixtureId)
	assert.Equal(t, f.Fixture.Id, a.FixtureId)
	assert.Equal(t, f.League.Id, h.LeagueId)
	assert.Equal(t, f.League.Id, a.LeagueId)
	assert.Equal(t, f.League.Season, h.Season)
	assert.Equal(t, f.League.Season, a.Season)
}
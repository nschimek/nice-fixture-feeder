package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetTLS(t *testing.T) {
	team := Team{Team: TeamTeam{Id: 100}}
	team.SetTLS(39, 2022)

	assert.Equal(t, TeamLeagueSeason{Id: TeamLeagueSeasonId{TeamId: 100, LeagueId: 39, Season: 2022}}, team.TeamLeagueSeason)
}

func TestGetTeamStatsIdCurrent(t *testing.T) {
	tls := TeamLeagueSeason{Id: TeamLeagueSeasonId{TeamId: 100, LeagueId: 39, Season: 2022}, MaxFixtureId: 101}

	a := tls.GetTeamStatsId(true)

	assert.Equal(t, &TeamStatsId{TeamId: 100, LeagueId: 39, FixtureId: 101, Season: 2022}, a)
}

func TestGetTeamStatsIdNext(t *testing.T) {
	tls := TeamLeagueSeason{Id: TeamLeagueSeasonId{TeamId: 100, LeagueId: 39, Season: 2022}, MaxFixtureId: 101}

	a := tls.GetTeamStatsId(false)

	assert.Equal(t, &TeamStatsId{TeamId: 100, LeagueId: 39, FixtureId: 0, Season: 2022, NextFixtureId: 101}, a)
}
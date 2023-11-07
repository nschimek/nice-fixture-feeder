package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTlsId(t *testing.T) {
	ts := TeamStats{Id: TeamStatsId{TeamId: 40, LeagueId: 39, FixtureId: 12345, Season: 2022}}

	a := ts.Id.GetTlsId()

	assert.Equal(t, TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, a)
}

func TestGetCurrentId(t *testing.T) {
	ts := TeamStats{Id: TeamStatsId{TeamId: 40, LeagueId: 39, FixtureId: 12345, Season: 2022, NextFixtureId: 12346}}

	a := ts.Id.GetCurrentId()

	assert.Equal(t, TeamStatsId{TeamId: 40, LeagueId: 39, FixtureId: 12345, Season: 2022, NextFixtureId: 0}, a)
}

func TestGetNextId(t *testing.T) {
	ts := TeamStats{Id: TeamStatsId{TeamId: 40, LeagueId: 39, FixtureId: 12345, Season: 2022, NextFixtureId: 12346}}

	a := ts.Id.GetNextId()

	assert.Equal(t, TeamStatsId{TeamId: 40, LeagueId: 39, FixtureId: 0, Season: 2022, NextFixtureId: 12346}, a)
}

func TestTeamStatsIncrement(t *testing.T) {
	s := &TeamStats{
		TeamStatsGoals: TeamStatsGoals{
			GoalsFor: TeamStatsHomeAwayTotal{Home: 10, Away: 7, Total: 17},
			GoalsAgainst: TeamStatsHomeAwayTotal{Home: 8, Away: 10, Total: 18},
		},
		TeamStatsFixtures: TeamStatsFixtures{
			FixturesPlayed: TeamStatsHomeAwayTotal{Home: 5, Away: 7, Total: 12},
			FixturesWins: TeamStatsHomeAwayTotal{Home: 3, Away: 1, Total: 4},
			FixturesLosses: TeamStatsHomeAwayTotal{Home: 1, Away: 2, Total: 3},
			FixturesDraws: TeamStatsHomeAwayTotal{Home: 2, Away: 3, Total: 5},
		},
	}

	s.TeamStatsGoals.GoalsFor.Increment(3, true)
	s.TeamStatsGoals.GoalsAgainst.Increment(2, true)
	s.TeamStatsFixtures.FixturesPlayed.Increment(1, true)
	s.TeamStatsFixtures.FixturesWins.Increment(1, true)

	s.TeamStatsGoals.GoalsFor.Increment(1, false)
	s.TeamStatsGoals.GoalsAgainst.Increment(5, false)
	s.TeamStatsFixtures.FixturesPlayed.Increment(1, false)
	s.TeamStatsFixtures.FixturesLosses.Increment(1, false)

	s.TeamStatsGoals.GoalsFor.Increment(3, false)
	s.TeamStatsGoals.GoalsAgainst.Increment(3, false)
	s.TeamStatsFixtures.FixturesPlayed.Increment(1, false)
	s.TeamStatsFixtures.FixturesDraws.Increment(1, false)

	// Goals For
	assert.Equal(t, 13, s.TeamStatsGoals.GoalsFor.Home)
	assert.Equal(t, 11, s.TeamStatsGoals.GoalsFor.Away)
	assert.Equal(t, 24, s.TeamStatsGoals.GoalsFor.Total)
	// Goals Against
	assert.Equal(t, 10, s.TeamStatsGoals.GoalsAgainst.Home)
	assert.Equal(t, 18, s.TeamStatsGoals.GoalsAgainst.Away)
	assert.Equal(t, 28, s.TeamStatsGoals.GoalsAgainst.Total)
	// Fixtures Played
	assert.Equal(t, 6, s.TeamStatsFixtures.FixturesPlayed.Home)
	assert.Equal(t, 9, s.TeamStatsFixtures.FixturesPlayed.Away)
	assert.Equal(t, 15, s.TeamStatsFixtures.FixturesPlayed.Total)
	// Fixtures Won
	assert.Equal(t, 4, s.TeamStatsFixtures.FixturesWins.Home)
	assert.Equal(t, 1, s.TeamStatsFixtures.FixturesWins.Away)
	assert.Equal(t, 5, s.TeamStatsFixtures.FixturesWins.Total)
	// Fixtures Lost
	assert.Equal(t, 1, s.TeamStatsFixtures.FixturesLosses.Home)
	assert.Equal(t, 3, s.TeamStatsFixtures.FixturesLosses.Away)
	assert.Equal(t, 4, s.TeamStatsFixtures.FixturesLosses.Total)
	// Fixtures Drawn
	assert.Equal(t, 2, s.TeamStatsFixtures.FixturesDraws.Home)
	assert.Equal(t, 4, s.TeamStatsFixtures.FixturesDraws.Away)
	assert.Equal(t, 6, s.TeamStatsFixtures.FixturesDraws.Total)
}
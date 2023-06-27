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
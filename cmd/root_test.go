package cmd

import (
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/request"
	"github.com/nschimek/nice-fixture-feeder/service"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCommandNoParams(t *testing.T) {
	ctx = new(commandContext)

	rootCmd.Run(&cobra.Command{}, []string{})

	ctx.startDate = time.Now().AddDate(0, 0, -1)
	ctx.endDate = time.Now()
}

func TestRootCommandVersion(t *testing.T) {
	var ok bool
	exitFunc = func(c int) {
		ok = true
	}

	err := rootCmd.PersistentPostRunE(&cobra.Command{Use: "version"}, []string{})

	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestRootSeason(t *testing.T) {
	var ok1, ok2 bool
	runSeasonFunc = func(lr request.LeagueRequest, tr request.TeamRequest) {
		ok1 = true
	}
	runFixturesFunc = func(fr request.FixtureRequest, ts service.TeamStatsService) {
		ok2 = true
	}
	services = &service.ServiceRegistry{TeamStats: &service.MockTeamStatsService{}}
	requests = &request.RequestRegistry{
		League: &request.MockLeagueRequest{}, 
		Team: &request.MockTeamRequest{}, 
		Fixture: &request.MockFixtureRequest{},
	}

	ctx = new(commandContext)
	ctx.season = true

	err := rootCmd.PersistentPostRunE(&cobra.Command{}, []string{})

	assert.Nil(t, err)
	assert.True(t, ok1)
	assert.True(t, ok2)
}

func TestRootNoSeason(t *testing.T) {
	var ok1, ok2 bool
	runSeasonFunc = func(lr request.LeagueRequest, tr request.TeamRequest) {
		ok1 = true
	}
	runFixturesFunc = func(fr request.FixtureRequest, ts service.TeamStatsService) {
		ok2 = true
	}
	services = &service.ServiceRegistry{TeamStats: &service.MockTeamStatsService{}}
	requests = &request.RequestRegistry{
		League: &request.MockLeagueRequest{}, 
		Team: &request.MockTeamRequest{}, 
		Fixture: &request.MockFixtureRequest{},
	}

	ctx = new(commandContext)

	err := rootCmd.PersistentPostRunE(&cobra.Command{}, []string{})

	assert.Nil(t, err)
	assert.False(t, ok1)
	assert.True(t, ok2)
}
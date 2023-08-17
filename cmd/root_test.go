package cmd

import (
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/request"
	req_mocks "github.com/nschimek/nice-fixture-feeder/request/mocks"
	"github.com/nschimek/nice-fixture-feeder/service"
	svc_mocks "github.com/nschimek/nice-fixture-feeder/service/mocks"
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
	runSeasonFunc = func(lr request.League, tr request.Team) {
		ok1 = true
	}
	runFixturesFunc = func(fr request.Fixture, ts service.TeamStats) {
		ok2 = true
	}
	services = &service.ServiceRegistry{TeamStats: &svc_mocks.TeamStats{}}
	requests = &request.RequestRegistry{
		League: &req_mocks.League{}, 
		Team: &req_mocks.Team{}, 
		Fixture: &req_mocks.Fixture{},
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
	runSeasonFunc = func(lr request.League, tr request.Team) {
		ok1 = true
	}
	runFixturesFunc = func(fr request.Fixture, ts service.TeamStats) {
		ok2 = true
	}
	services = &service.ServiceRegistry{TeamStats: &svc_mocks.TeamStats{}}
	requests = &request.RequestRegistry{
		League: &req_mocks.League{}, 
		Team: &req_mocks.Team{}, 
		Fixture: &req_mocks.Fixture{},
	}

	ctx = new(commandContext)

	err := rootCmd.PersistentPostRunE(&cobra.Command{}, []string{})

	assert.Nil(t, err)
	assert.False(t, ok1)
	assert.True(t, ok2)
}
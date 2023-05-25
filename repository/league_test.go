package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type LeagueRepositoryTestSuite struct {
	suite.Suite
	seasons []model.LeagueSeason
	league model.LeagueLeague
	leagues []model.League
	mockDatabase *mocks.Database
}

func TestLeagueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LeagueRepositoryTestSuite))
}

func (suite *LeagueRepositoryTestSuite) SetupTest() {
	suite.seasons = []model.LeagueSeason{{Season: 2022}}
	suite.league = model.LeagueLeague{Id: 1}
	suite.leagues = []model.League{{League: suite.league, Seasons: suite.seasons}}
	suite.mockDatabase = &mocks.Database{}
}

func (suite *LeagueRepositoryTestSuite) TestUpsertLeaguesSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	suite.mockDatabase.On("Upsert", &suite.leagues[0]).Return(r)

	repo := &LeagueRepository{DB: suite.mockDatabase}
	actual := repo.Upsert(suite.leagues)

	suite.Equal(0, actual.Error["league"])
	suite.Equal(0, actual.Error["season"])
	suite.Equal(1, actual.Success["league"])
	suite.Equal(1, actual.Success["season"])
}

func (suite *LeagueRepositoryTestSuite) TestUpsertLeagueError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	suite.mockDatabase.On("Upsert", &suite.leagues[0]).Return(r)

	repo := &LeagueRepository{DB: suite.mockDatabase}
	actual := repo.Upsert(suite.leagues)

	suite.Equal(1, actual.Error["league"])
	suite.Equal(1, actual.Error["season"])
	suite.Equal(0, actual.Success["league"])
	suite.Equal(0, actual.Success["season"])
}

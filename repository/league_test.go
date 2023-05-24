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
}

func (suite *LeagueRepositoryTestSuite) SetupTest() {
	suite.seasons = []model.LeagueSeason{{Season: 2022}}
	suite.league = model.LeagueLeague{Id: 1}
	suite.leagues = []model.League{{League: suite.league, Seasons: suite.seasons}}
}

func (suite *LeagueRepositoryTestSuite) TestUpsertLeaguesSuccess() {
	r1 := core.DatabaseResult{RowsAffected: 1, Error: nil}
	r2 := core.DatabaseResult{RowsAffected: 2, Error: nil}

	mockDatabase := &mocks.Database{}
	mockDatabase.On("UpsertWithOmit", &suite.leagues[0], "Seasons").Return(r1)
	mockDatabase.On("Upsert", suite.leagues[0].Seasons).Return(r2)

	repo := &LeagueRepository{DB: mockDatabase}
	actual := repo.UpsertLeagues(suite.leagues)

	suite.Equal(0, actual.Error["league"])
	suite.Equal(0, actual.Error["season"])
	suite.Equal(1, actual.Success["league"])
	suite.Equal(1, actual.Success["season"])
}

func (suite *LeagueRepositoryTestSuite) TestUpsertLeagueErrorLeague() {
	r1 := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	mockDatabase := &mocks.Database{}
	mockDatabase.On("UpsertWithOmit", &suite.leagues[0], "Seasons").Return(r1)

	repo := &LeagueRepository{DB: mockDatabase}
	actual := repo.UpsertLeagues(suite.leagues)

	suite.Equal(1, actual.Error["league"])
	suite.Equal(0, actual.Success["league"])
	suite.True(suite.leagues[0].ModelError.HasErrors())
}

func (suite *LeagueRepositoryTestSuite) TestUpsertLeagueErrorSeason() {
	r1 := core.DatabaseResult{RowsAffected: 1, Error: nil}
	r2 := core.DatabaseResult{RowsAffected: 2, Error: errors.New("test")}

	mockDatabase := &mocks.Database{}
	mockDatabase.On("UpsertWithOmit", &suite.leagues[0], "Seasons").Return(r1)
	mockDatabase.On("Upsert", suite.leagues[0].Seasons).Return(r2)

	repo := &LeagueRepository{DB: mockDatabase}
	actual := repo.UpsertLeagues(suite.leagues)

	suite.Equal(0, actual.Error["league"])
	suite.Equal(1, actual.Error["season"])
	suite.Equal(1, actual.Success["league"])
	suite.Equal(0, actual.Success["season"])
}

func TestLeagueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(LeagueRepositoryTestSuite))
}
package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type leagueRepositoryTestSuite struct {
	suite.Suite
	seasons []model.LeagueSeason
	league model.LeagueLeague
	leagues []model.League
	mockDatabase *core.MockDatabase
}

func TestLeagueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(leagueRepositoryTestSuite))
}

func (s *leagueRepositoryTestSuite) SetupTest() {
	s.seasons = []model.LeagueSeason{{Season: 2022}}
	s.league = model.LeagueLeague{Id: 1}
	s.leagues = []model.League{{League: s.league, Seasons: s.seasons}}
	s.mockDatabase = &core.MockDatabase{}
}

func (s *leagueRepositoryTestSuite) TestUpsertLeaguesSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.leagues).Return(r)

	repo := NewLeagueRepository(s.mockDatabase)
	actual, err := repo.Upsert(s.leagues)

	s.Equal(s.leagues, actual)
	s.Nil(err)
}

func (s *leagueRepositoryTestSuite) TestUpsertLeagueError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.leagues).Return(r)

	repo := NewLeagueRepository(s.mockDatabase)
	actual, err := repo.Upsert(s.leagues)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

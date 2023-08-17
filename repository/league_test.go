package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type leagueRepositoryTestSuite struct {
	suite.Suite
	seasons []model.LeagueSeason
	league model.LeagueLeague
	leagues []model.League
	mockDatabase *mocks.Database
	repo *League
}

func TestLeagueRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(leagueRepositoryTestSuite))
}

func (s *leagueRepositoryTestSuite) SetupTest() {
	s.seasons = []model.LeagueSeason{{Season: 2022}}
	s.league = model.LeagueLeague{Id: 1}
	s.leagues = []model.League{{League: s.league, Seasons: s.seasons}}
	s.mockDatabase = &mocks.Database{}
	s.repo = NewLeague(s.mockDatabase)
}

func (s *leagueRepositoryTestSuite) TestUpsertLeaguesSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.leagues).Return(r)

	actual, err := s.repo.Upsert(s.leagues)

	s.Equal(s.leagues, actual)
	s.Nil(err)
}

func (s *leagueRepositoryTestSuite) TestUpsertLeagueError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.leagues).Return(r)

	actual, err := s.repo.Upsert(s.leagues)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

func (s *leagueRepositoryTestSuite) TestUpsertEmptyAndNil() {
	a1, e1 := s.repo.Upsert([]model.League{})
	s.Nil(a1)
	s.Nil(e1)
	a2, e2 := s.repo.Upsert(nil)
	s.Nil(a2)
	s.Nil(e2)
}
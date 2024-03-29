package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type teamStatsRepositoryTestSuite struct {
	suite.Suite
	mockDatabase *mocks.Database
	repo TeamStats
	teamStats []model.TeamStats
}

func TestTeamStatsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(teamStatsRepositoryTestSuite))
}

func (s *teamStatsRepositoryTestSuite) SetupTest() {
	s.mockDatabase = &mocks.Database{}
	s.repo = NewTeamStats(s.mockDatabase)
	s.teamStats = []model.TeamStats{
		{Id: model.TeamStatsId{TeamId: 1, LeagueId: 2, Season: 2022, FixtureId: 100}},
		{Id: model.TeamStatsId{TeamId: 1, LeagueId: 2, Season: 2022, FixtureId: 101}},
		{Id: model.TeamStatsId{TeamId: 1, LeagueId: 2, Season: 2022, FixtureId: 102}},
	}
}

func (s *teamStatsRepositoryTestSuite) TestUpsertSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.teamStats).Return(r)

	actual, err := s.repo.Upsert(s.teamStats)

	s.Equal(s.teamStats, actual)
	s.Nil(err)
}

func (s *teamStatsRepositoryTestSuite) TestUpsertError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.teamStats).Return(r)

	actual, err := s.repo.Upsert(s.teamStats)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

func (s *teamStatsRepositoryTestSuite) TestGetByIdFound() {
	var entity model.TeamStats
	id := s.teamStats[0].Id

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 1})

	resp, err := s.repo.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Equal(&entity, resp) // due to mocking, expect the response to just be the entity passed through
	s.Nil(err)
}

func (s *teamStatsRepositoryTestSuite) TestGetByIdNotFound() {
	var entity model.TeamStats
	id := model.TeamStatsId{TeamId: 99, LeagueId: 9, Season: 2022}

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 0})

	resp, err := s.repo.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Nil(resp)
	s.Nil(err)
}

func (s *teamStatsRepositoryTestSuite) TestGetByIdError() {
	var entity model.TeamStats
	id := model.TeamStatsId{TeamId: 99, LeagueId: 9, Season: 2022, FixtureId: 9999}

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")})

	resp, err := s.repo.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Nil(resp)
	s.ErrorContains(err, "test")
}

func (s *teamStatsRepositoryTestSuite) TestUpsertEmptyAndNil() {
	a1, e1 := s.repo.Upsert([]model.TeamStats{})
	s.Nil(a1)
	s.Nil(e1)
	a2, e2 := s.repo.Upsert(nil)
	s.Nil(a2)
	s.Nil(e2)
}

func (s *teamStatsRepositoryTestSuite) TestUpsertOne() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}
	exp := s.teamStats[0:1] // have to create the slice first, then can get the address in the mock

	s.mockDatabase.EXPECT().Upsert(&exp).Return(r)

	actual, err := s.repo.UpsertOne(s.teamStats[0])

	s.Equal(s.teamStats[0], actual)
	s.Nil(err)
}

func (s *teamStatsRepositoryTestSuite) TestUpsertOneError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}
	exp := s.teamStats[0:1] // have to create the slice first, then can get the address in the mock

	s.mockDatabase.EXPECT().Upsert(&exp).Return(r)

	actual, err := s.repo.UpsertOne(s.teamStats[0])

	s.Equal(model.TeamStats{}, actual)
	s.ErrorContains(err, "test")
}
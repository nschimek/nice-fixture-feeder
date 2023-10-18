package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type teamLeagueSeasonRepositoryTestSuite struct {
	suite.Suite
	mockDatabase *mocks.Database
	repo TeamLeagueSeason
	teamLeagueSeasons []model.TeamLeagueSeason
}

func TestTeamLeagueSeasonRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(teamLeagueSeasonRepositoryTestSuite))
}

func (s *teamLeagueSeasonRepositoryTestSuite) SetupTest() {
	s.mockDatabase = &mocks.Database{}
	s.repo = NewTeamLeagueSeason(s.mockDatabase)
	s.teamLeagueSeasons = []model.TeamLeagueSeason{
		{Id: model.TeamLeagueSeasonId{TeamId: 42, LeagueId: 39, Season: 2022}},
		{Id: model.TeamLeagueSeasonId{TeamId: 43, LeagueId: 39, Season: 2022}},
		{Id: model.TeamLeagueSeasonId{TeamId: 44, LeagueId: 39, Season: 2022}},
	}
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestUpsertSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.teamLeagueSeasons).Return(r)

	actual, err := s.repo.Upsert(s.teamLeagueSeasons)

	s.Equal(s.teamLeagueSeasons, actual)
	s.Nil(err)
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestUpsertError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.teamLeagueSeasons).Return(r)

	actual, err := s.repo.Upsert(s.teamLeagueSeasons)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestGetByIdFound() {
	var entity model.TeamLeagueSeason
	id := s.teamLeagueSeasons[0].Id

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 1})

	resp, err := s.repo.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Equal(&entity, resp) // due to mocking, expect the response to just be the entity passed through
	s.Nil(err)
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestGetByIdNotFound() {
	var entity model.TeamLeagueSeason
	id := model.TeamLeagueSeasonId{TeamId: 99, LeagueId: 9, Season: 2022}

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 0})

	resp, err := s.repo.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Nil(resp)
	s.Nil(err)
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestGetByIdError() {
	var entity model.TeamLeagueSeason
	id := model.TeamLeagueSeasonId{TeamId: 99, LeagueId: 9, Season: 2022}

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")})

	resp, err := s.repo.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Nil(resp)
	s.ErrorContains(err, "test")
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestUpsertEmptyAndNil() {
	a1, e1 := s.repo.Upsert([]model.TeamLeagueSeason{})
	s.Nil(a1)
	s.Nil(e1)
	a2, e2 := s.repo.Upsert(nil)
	s.Nil(a2)
	s.Nil(e2)
}

package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type teamLeagueSeasonRepositoryTestSuite struct {
	suite.Suite
	mockDatabase *core.MockDatabase
	teamLeagueSeasonRepository TeamLeagueSeasonRepository
}

func TestTeamLeagueSeasonRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(teamLeagueSeasonRepositoryTestSuite))
}

func (s *teamLeagueSeasonRepositoryTestSuite) SetupTest() {
	s.mockDatabase = &core.MockDatabase{}
	s.teamLeagueSeasonRepository = NewTeamLeagueSeasonRepository(s.mockDatabase)
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestGetByIdFound() {
	var entity model.TeamLeagueSeason
	id := model.TeamLeagueSeason{TeamId: 42, LeagueId: 39, Season: 2022}

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 1})

	resp := s.teamLeagueSeasonRepository.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Equal(&entity, resp) // due to mocking, expect the response to just be the entity passed through
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestGetByIdNotFound() {
	var entity model.TeamLeagueSeason
	id := model.TeamLeagueSeason{TeamId: 99, LeagueId: 9, Season: 2022}

	s.mockDatabase.EXPECT().GetById(id, &entity).Return(core.DatabaseResult{RowsAffected: 0})

	resp := s.teamLeagueSeasonRepository.GetById(id)

	s.mockDatabase.AssertCalled(s.T(), "GetById", id, &entity)
	s.Nil(resp)
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestSaveSuccess() {
	e := model.TeamLeagueSeason{TeamId: 42, LeagueId: 39, Season: 2022}
	
	s.mockDatabase.EXPECT().Save(&e).Return(core.DatabaseResult{RowsAffected: 1})

	a, err := s.teamLeagueSeasonRepository.Save(&e)

	s.mockDatabase.AssertCalled(s.T(), "Save", &e)
	s.Equal(&e, a)
	s.Nil(err)
}

func (s *teamLeagueSeasonRepositoryTestSuite) TestSaveError() {
	e := model.TeamLeagueSeason{TeamId: 42, LeagueId: 39, Season: 2022}
	
	s.mockDatabase.EXPECT().Save(&e).Return(core.DatabaseResult{Error: errors.New("test")})

	a, err := s.teamLeagueSeasonRepository.Save(&e)

	s.mockDatabase.AssertCalled(s.T(), "Save", &e)
	s.Nil(a)
	s.ErrorContains(err, "test")
}
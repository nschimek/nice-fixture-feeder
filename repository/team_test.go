package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type teamRepositoryTestSuite struct {
	suite.Suite
	teams []model.Team
	mockDatabase *core.MockDatabase
}

func TestTeamRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(teamRepositoryTestSuite))
}

func (s *teamRepositoryTestSuite) SetupTest() {
	s.teams = []model.Team{
		{Team: model.TeamTeam{Id: 1, Name: "Test"}},
		{Team: model.TeamTeam{Id: 2, Name: "Test"}},
	}
	s.mockDatabase = &core.MockDatabase{}
}

func (s *teamRepositoryTestSuite) TestUpsertTeamsSuccess() {
	r := core.DatabaseResult{RowsAffected: 2}

	s.mockDatabase.EXPECT().Upsert(&s.teams).Return(r)

	repo := NewTeamRepository(s.mockDatabase)
	actual, err := repo.Upsert(s.teams)

	s.Equal(s.teams, actual)
	s.Nil(err)
}

func (s *teamRepositoryTestSuite) TestUpsertTeamError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.teams).Return(r)

	repo := NewTeamRepository(s.mockDatabase)
	actual, err := repo.Upsert(s.teams)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}
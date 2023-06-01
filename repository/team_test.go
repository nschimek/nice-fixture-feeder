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
		{
			Team: model.TeamTeam{
				Id: 1,
				Name: "Test",
			},
		},
		{
			Team: model.TeamTeam{
				Id: 2,
				Name: "Test",
			},
		},
	}
	s.mockDatabase = &core.MockDatabase{}
}

func (s *teamRepositoryTestSuite) TestUpsertteamsSuccess() {
	r := core.DatabaseResult{RowsAffected: 2, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.teams).Return(r)

	repo := &TeamRepository{DB: s.mockDatabase}
	actual := repo.Upsert(s.teams)

	s.Equal(0, actual.Error["team"])
	s.Equal(0, actual.Error["team_league_season"])
	s.Equal(2, actual.Success["team"])
	s.Equal(2, actual.Success["team_league_season"])
}

func (s *teamRepositoryTestSuite) TestUpsertteamError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.teams).Return(r)

	repo := &TeamRepository{DB: s.mockDatabase}
	actual := repo.Upsert(s.teams)

	s.Equal(2, actual.Error["team"])
	s.Equal(2, actual.Error["team_league_season"])
	s.Equal(0, actual.Success["team"])
	s.Equal(0, actual.Success["team_league_season"])
}
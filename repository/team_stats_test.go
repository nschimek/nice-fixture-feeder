package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type teamStatsRepositoryTestSuite struct {
	suite.Suite
	mockDatabase *core.MockDatabase
	teamStatsRepostiory TeamStatsRepository
	teamStats []model.TeamStats
}

func TestTeamStatsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(teamStatsRepositoryTestSuite))
}

func (s *teamStatsRepositoryTestSuite) SetupTest() {
	s.mockDatabase = &core.MockDatabase{}
	s.teamStatsRepostiory = NewTeamStatsRepository(s.mockDatabase)
	s.teamStats = []model.TeamStats{
		{TeamStatsId: model.TeamStatsId{TeamId: 1, LeagueId: 2, Season: 2022, FixtureId: 100}},
		{TeamStatsId: model.TeamStatsId{TeamId: 1, LeagueId: 2, Season: 2022, FixtureId: 101}},
		{TeamStatsId: model.TeamStatsId{TeamId: 1, LeagueId: 2, Season: 2022, FixtureId: 102}},
	}
}

func (s *teamStatsRepositoryTestSuite) TestUpsertFixturesSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.teamStats).Return(r)

	actual, err := s.teamStatsRepostiory.Upsert(s.teamStats)

	s.Equal(s.teamStats, actual)
	s.Nil(err)
}

func (s *teamStatsRepositoryTestSuite) TestUpsertFixtureError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.teamStats).Return(r)

	actual, err := s.teamStatsRepostiory.Upsert(s.teamStats)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

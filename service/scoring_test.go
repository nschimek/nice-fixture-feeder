package service

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/nschimek/nice-fixture-feeder/service/mocks"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
	score_mocks "github.com/nschimek/nice-fixture-feeder/service/scores/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type scoringTestSuite struct {
	suite.Suite
	mockFixtureRepo *repo_mocks.Fixture
	mockFixtureScoreRepo *repo_mocks.UpsertRepository[model.FixtureScore]
	mockStatsService *mocks.TeamStats
	mockStatusService *mocks.FixtureStatus
	mockScorePS *score_mocks.PointsStrength
	scoring scoring
	fixtures []model.Fixture
	fixtureScores []model.FixtureScore
}

func TestScoringTestSuite(t *testing.T) {
	suite.Run(t, new(scoringTestSuite))
}

func (s *scoringTestSuite) SetupTest() {
	s.mockFixtureRepo = new(repo_mocks.Fixture)
	s.mockFixtureScoreRepo = new(repo_mocks.UpsertRepository[model.FixtureScore])
	s.mockStatsService = new(mocks.TeamStats)
	s.mockStatusService = new(mocks.FixtureStatus)
	s.mockScorePS = new(score_mocks.PointsStrength)
	s.mockScorePS.EXPECT().SetStatsFunc(mock.AnythingOfType("func(*model.Fixture) (*model.TeamStats, *model.TeamStats, error)"))
	s.scoring = *NewScoring(
		&scores.ScoreRegistry{
			AllScores: []scores.Score{s.mockScorePS},
			PointsStrength: s.mockScorePS,
		},
		s.mockFixtureRepo,
		s.mockFixtureScoreRepo,
		s.mockStatsService,
		s.mockStatusService,
	)
	s.fixtures = []model.Fixture{
		{Fixture: model.FixtureFixture{Id: 100}},
		{Fixture: model.FixtureFixture{Id: 101}},
		{Fixture: model.FixtureFixture{Id: 102}},
	}
	s.fixtureScores = []model.FixtureScore{
		{FixtureId: 100, ScoreId: 1, Value: 67},
		{FixtureId: 101, ScoreId: 1, Value: 51},
	}
}

func (s *scoringTestSuite) TestAddFixturesFromMinMap() {
	fixturesMap := map[int]model.Fixture{
		100: s.fixtures[0],
		101: s.fixtures[1],
	}
	notIn := []int{100, 101}
	fixturesMinMap := map[model.TeamLeagueSeasonId]int{
		{TeamId: 1, LeagueId: 10, Season: 2022}: 99,
		{TeamId: 2, LeagueId: 10, Season: 2022}: 99,
	}

	s.mockFixtureRepo.EXPECT().
		GetFutureFixturesByTLS(model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 10, Season: 2022}, 99, notIn).
		Return([]model.Fixture{{Fixture: model.FixtureFixture{Id: 102}}}, nil)
	s.mockFixtureRepo.EXPECT().
		GetFutureFixturesByTLS(model.TeamLeagueSeasonId{TeamId: 2, LeagueId: 10, Season: 2022}, 99, notIn).
		Return([]model.Fixture{{Fixture: model.FixtureFixture{Id: 103}}}, nil)

	// set fixtures map directly to also test this method
	s.scoring.SetFixtures(fixturesMap)
	s.scoring.AddFixturesFromMinMap(fixturesMinMap)

	s.Len(s.scoring.fixturesMap, 4)
	s.Contains(s.scoring.fixturesMap, 100)
	s.Contains(s.scoring.fixturesMap, 101)
	s.Contains(s.scoring.fixturesMap, 102)
	s.Contains(s.scoring.fixturesMap, 103)
}

func (s *scoringTestSuite) TestScore() {
	s.scoring.SetFixtures(map[int]model.Fixture{
		100: s.fixtures[0],
		101: s.fixtures[1],
		102: s.fixtures[2],
	})

	s.mockScorePS.EXPECT().CanScore(&s.fixtures[0]).Return(true)
	s.mockScorePS.EXPECT().CanScore(&s.fixtures[1]).Return(false)
	s.mockScorePS.EXPECT().CanScore(&s.fixtures[2]).Return(true)
	s.mockScorePS.EXPECT().Score(&s.fixtures[0]).Return(&s.fixtureScores[0], nil)
	s.mockScorePS.EXPECT().Score(&s.fixtures[2]).Return(nil, errors.New("test")) // coverage
	s.mockFixtureScoreRepo.EXPECT().Upsert(s.fixtureScores[0:1]).Return(s.fixtureScores[0:1], nil)
	
	s.scoring.Score()
	
	s.mockScorePS.AssertExpectations(s.T())
	s.mockFixtureScoreRepo.AssertCalled(s.T(), "Upsert", s.fixtureScores[0:1])
}

func (s *scoringTestSuite) TestGetStatsTeam() {
	fixture := model.Fixture{
		League: model.FixtureLeague{Id: 10, Season: 2022, Round: 5}, 
		Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 1}}, 
		Fixture: model.FixtureFixture{Id: 100, Status: model.FixtureStatusId{Id: "FT"}},
	}
	stats := model.TeamStats{Id: model.TeamStatsId{FixtureId: 100}}

	s.mockStatusService.EXPECT().IsScheduled("FT").Return(false)
	s.mockStatsService.EXPECT().GetById(model.TeamStatsId{
		TeamId: 1, LeagueId: 10, Season: 2022, NextFixtureId: 100}, false).Return(&stats, nil)

	res, err := s.scoring.getStatsTeam(&fixture, true)

	s.NoError(err)
	s.Equal(&stats, res)
}

func (s *scoringTestSuite) TestGetStatsTeamScheduled() {
	fixture := model.Fixture{
		League: model.FixtureLeague{Id: 10, Season: 2022, Round: 5}, 
		Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 1}}, 
		Fixture: model.FixtureFixture{Id: 100, Status: model.FixtureStatusId{Id: "SC"}},
	}
	stats := model.TeamStats{Id: model.TeamStatsId{FixtureId: 100}}

	s.mockStatusService.EXPECT().IsScheduled("SC").Return(true)
	s.mockStatsService.EXPECT().GetByIdWithTLS(model.TeamStatsId{
		TeamId: 1, LeagueId: 10, Season: 2022, NextFixtureId: 100}, false).Return(&stats, nil)

	res, err := s.scoring.getStatsTeam(&fixture, true)

	s.NoError(err)
	s.Equal(&stats, res)
}

func (s *scoringTestSuite) TestGetStatsTeamEarlyRound() {
	fixture := model.Fixture{
		League: model.FixtureLeague{Id: 10, Season: 2022, Round: 1}, 
		Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 1}}, 
		Fixture: model.FixtureFixture{Id: 100, Status: model.FixtureStatusId{Id: "FT"}},
	}
	stats := model.TeamStats{Id: model.TeamStatsId{FixtureId: 100}}

	s.mockStatusService.EXPECT().IsScheduled("FT").Return(false)
	s.mockStatsService.EXPECT().GetByIdWithTLS(model.TeamStatsId{
		TeamId: 1, LeagueId: 10, Season: 2021, NextFixtureId: 100}, true).Return(&stats, nil)

	res, err := s.scoring.getStatsTeam(&fixture, true)

	s.NoError(err)
	s.Equal(&stats, res)
}
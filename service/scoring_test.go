package service

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/nschimek/nice-fixture-feeder/service/mocks"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
	score_mocks "github.com/nschimek/nice-fixture-feeder/service/scores/mocks"

	"github.com/stretchr/testify/suite"
)

type scoringTestSuite struct {
	suite.Suite
	mockFixtureRepo *repo_mocks.Fixture
	mockFixtureScoreRepo *repo_mocks.UpsertRepository[model.FixtureScore]
	mockStatsService *mocks.TeamStats
	mockStatusService *mocks.FixtureStatus
	mockScorePS *score_mocks.PointsStrength
	scoring Scoring
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
	s.scoring = NewScoring(
		&scores.ScoreRegistry{
			AllScores: []scores.Score{s.mockScorePS},
			PointsStrength: s.mockScorePS,
		},
		s.mockFixtureRepo,
		s.mockFixtureScoreRepo,
		s.mockStatsService,
		s.mockStatusService,
	)
}
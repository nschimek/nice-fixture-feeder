package service

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type teamLeagueSeasonServiceTestSuite struct {
	suite.Suite
	mockRepo *repo_mocks.TeamLeagueSeason
	tlsService *teamLeagueSeason
	tls []model.TeamLeagueSeason
}

func TestTeamLeagueSeasonServiceTestSuite(t *testing.T) {
	suite.Run(t, new(teamLeagueSeasonServiceTestSuite))
}

func (s *teamLeagueSeasonServiceTestSuite) SetupTest() {
	s.mockRepo = &repo_mocks.TeamLeagueSeason{}
	s.tlsService = NewTeamLeagueSeason(s.mockRepo)
	s.tls = []model.TeamLeagueSeason{
		{Id: model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 2, Season: 2022}, MaxFixtureId: 100},
		{Id: model.TeamLeagueSeasonId{TeamId: 2, LeagueId: 2, Season: 2022}, MaxFixtureId: 101},
		{Id: model.TeamLeagueSeasonId{TeamId: 3, LeagueId: 2, Season: 2022}, MaxFixtureId: 102},
	}
}

func (s *teamLeagueSeasonServiceTestSuite) TestGetByIdFound() {
	id := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 2, Season: 2022}

	s.mockRepo.EXPECT().GetById(id).Return(&s.tls[0], nil)

	tls, err := s.tlsService.GetById(id)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.tls[0], *tls)
	assert.Contains(s.T(), s.tlsService.tlsMap, id)
}
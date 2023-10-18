package service

import (
	"errors"
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

func (s* teamLeagueSeasonServiceTestSuite) TestGetByIdMapFound() {
	id := s.tls[0].Id
	s.tlsService.AddToMap(&s.tls[0]) // also tests AddToMap for us

	tls, err := s.tlsService.GetById(id)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.tls[0], *tls)
	s.mockRepo.AssertNotCalled(s.T(), "GetById", id)
	assert.Contains(s.T(), s.tlsService.tlsMap, id) 
}

func (s *teamLeagueSeasonServiceTestSuite) TestGetByIdRepoFound() {
	id := s.tls[0].Id

	s.mockRepo.EXPECT().GetById(id).Return(&s.tls[0], nil)

	tls, err := s.tlsService.GetById(id)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.tls[0], *tls)
	assert.Contains(s.T(), s.tlsService.tlsMap, id)
}

func (s *teamLeagueSeasonServiceTestSuite) TestGetByIdNotFound() {
	id := s.tls[0].Id
	
	s.mockRepo.EXPECT().GetById(id).Return(nil, errors.New("test"))

	tls, err := s.tlsService.GetById(id)

	assert.NotNil(s.T(), err)
	assert.ErrorContains(s.T(), err, "could not get TLS")
	assert.Nil(s.T(), tls)
}

func (s *teamLeagueSeasonServiceTestSuite) TestPersist() {
	s.tlsService.AddToMap(&s.tls[0])

	s.mockRepo.EXPECT().Upsert(s.tls[0:1]).Return(s.tls[0:1], nil)

	s.tlsService.Persist()

	s.mockRepo.AssertCalled(s.T(), "Upsert", s.tls[0:1])
}
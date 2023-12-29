package service

import (
	"errors"
	"testing"

	core_mocks "github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type teamLeagueSeasonServiceTestSuite struct {
	suite.Suite
	mockRepo *repo_mocks.TeamLeagueSeason
	mockCache *core_mocks.Cache[model.TeamLeagueSeason]
	tlsService *teamLeagueSeason
	tls []model.TeamLeagueSeason
}

func TestTeamLeagueSeasonServiceTestSuite(t *testing.T) {
	suite.Run(t, new(teamLeagueSeasonServiceTestSuite))
}

func (s *teamLeagueSeasonServiceTestSuite) SetupTest() {
	s.mockRepo = &repo_mocks.TeamLeagueSeason{}	
	s.mockCache = &core_mocks.Cache[model.TeamLeagueSeason]{}
	s.tlsService = NewTeamLeagueSeason(s.mockRepo, s.mockCache)
	s.tls = []model.TeamLeagueSeason{
		{Id: model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 2, Season: 2022}, MaxFixtureId: 100},
		{Id: model.TeamLeagueSeasonId{TeamId: 2, LeagueId: 2, Season: 2022}, MaxFixtureId: 101},
		{Id: model.TeamLeagueSeasonId{TeamId: 3, LeagueId: 2, Season: 2022}, MaxFixtureId: 102},
	}
}

func (s* teamLeagueSeasonServiceTestSuite) TestGetByIdCacheHit() {
	id := s.tls[0].Id

	s.mockCache.EXPECT().Get(id).Return(&s.tls[0], nil)

	tls, err := s.tlsService.GetById(id)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.tls[0], *tls)
	s.mockRepo.AssertNotCalled(s.T(), "GetById", id) 
}

func (s *teamLeagueSeasonServiceTestSuite) TestGetByIdRepoFound() {
	id := s.tls[0].Id

	s.mockCache.EXPECT().Get(id).Return(nil, nil)
	s.mockRepo.EXPECT().GetById(id).Return(&s.tls[0], nil)
	s.mockCache.EXPECT().Set(id, &s.tls[0]).Return(nil)

	tls, err := s.tlsService.GetById(id)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), s.tls[0], *tls)
}

func (s *teamLeagueSeasonServiceTestSuite) TestGetByIdNotFound() {
	id := s.tls[0].Id
	
	s.mockCache.EXPECT().Get(id).Return(nil, nil)
	s.mockRepo.EXPECT().GetById(id).Return(nil, errors.New("test"))

	tls, err := s.tlsService.GetById(id)

	assert.NotNil(s.T(), err)
	assert.ErrorContains(s.T(), err, "could not get TLS")
	assert.Nil(s.T(), tls)
}

func (s *teamLeagueSeasonServiceTestSuite) TestPersistOne() {

	s.mockRepo.EXPECT().UpsertOne(s.tls[0]).Return(&s.tls[0], nil)
	s.mockCache.EXPECT().Set(s.tls[0].Id, &s.tls[0]).Return(nil)

	s.tlsService.PersistOne(&s.tls[0])

	s.mockRepo.AssertExpectations(s.T())
	s.mockCache.AssertExpectations(s.T())
}
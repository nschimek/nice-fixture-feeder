package service

import (
	"testing"

	core_mocks "github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type fixtureStatusServiceTestSuite struct {
	suite.Suite
	mockRepository *mocks.FixtureStatus
	mockCache *core_mocks.Cache[model.FixtureStatus]
	fixtureStatusService FixtureStatus
	statuses []model.FixtureStatus
}

func TestFixtureStatusServiceTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureStatusServiceTestSuite))
}

func (s *fixtureStatusServiceTestSuite) SetupTest() {
	s.mockRepository = &mocks.FixtureStatus{}
	s.mockCache = &core_mocks.Cache[model.FixtureStatus]{}
	s.statuses = []model.FixtureStatus{
		{Id: "FT", Type: "FI"},
		{Id: "NS", Type: "SC"},
		{Id: "LIVE", Type: "IP"},
	}
	s.mockRepository.EXPECT().GetAll().Return(s.statuses, nil)
	s.mockCache.EXPECT().Set("FT", &s.statuses[0]).Return(nil)
	s.mockCache.EXPECT().Set("NS", &s.statuses[1]).Return(nil)
	s.mockCache.EXPECT().Set("LIVE", &s.statuses[2]).Return(nil)
	s.fixtureStatusService = NewFixtureStatus(s.mockRepository, s.mockCache)
}

func (s *fixtureStatusServiceTestSuite) TestGetType() {
	s.mockCache.EXPECT().Get("FT").Return(&s.statuses[0], nil)

	s.Equal("FI", s.fixtureStatusService.GetType("ft"))
}

func (s *fixtureStatusServiceTestSuite) TestFinished() {
	s.mockCache.EXPECT().Get("FT").Return(&s.statuses[0], nil)
	s.mockCache.EXPECT().Get("NS").Return(&s.statuses[1], nil)


	s.True(s.fixtureStatusService.IsFinished("FT"))
	s.False(s.fixtureStatusService.IsFinished("NS"))
}

func (s *fixtureStatusServiceTestSuite) TestScheduled() {
	s.mockCache.EXPECT().Get("NS").Return(&s.statuses[1], nil)
	s.mockCache.EXPECT().Get("LIVE").Return(&s.statuses[2], nil)

	s.True(s.fixtureStatusService.IsScheduled("NS"))
	s.False(s.fixtureStatusService.IsScheduled("LIVE"))
}
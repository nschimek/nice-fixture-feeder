package service

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/stretchr/testify/suite"
)

type fixtureStatusServiceTestSuite struct {
	suite.Suite
	mockRepository *repository.MockFixtureStatusRepository
	fixtureStatusService FixtureStatusService
	statuses []model.FixtureStatus
}

func TestFixtureStatusServiceTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureStatusServiceTestSuite))
}

func (s *fixtureStatusServiceTestSuite) SetupTest() {
	s.mockRepository = &repository.MockFixtureStatusRepository{}
	s.fixtureStatusService = NewFixtureStatusService(s.mockRepository)
	s.statuses = []model.FixtureStatus{
		{Id: "FT", Type: "FI"},
		{Id: "NS", Type: "SC"},
		{Id: "LIVE", Type: "IP"},
	}
}

func (s *fixtureStatusServiceTestSuite) TestGetMap() {
	s.mockRepository.EXPECT().GetAll().Return(s.statuses)

	s.Equal("FI", s.fixtureStatusService.GetType("ft"))
	s.Equal("SC", s.fixtureStatusService.GetType("NS"))
	s.Equal("IP", s.fixtureStatusService.GetType("live"))

	// we should only call the repository once
	s.mockRepository.AssertNumberOfCalls(s.T(), "GetAll", 1)
}

func (s *fixtureStatusServiceTestSuite) TestFinished() {
	s.mockRepository.EXPECT().GetAll().Return(s.statuses)

	s.True(s.fixtureStatusService.IsFinished("FT"))
	s.False(s.fixtureStatusService.IsFinished("NS"))
}

func (s *fixtureStatusServiceTestSuite) TestScheduled() {
	s.mockRepository.EXPECT().GetAll().Return(s.statuses)

	s.True(s.fixtureStatusService.IsScheduled("NS"))
	s.False(s.fixtureStatusService.IsScheduled("LIVE"))
}
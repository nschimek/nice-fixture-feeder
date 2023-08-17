package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type fixtureStatusRepositoryTestSuite struct {
	suite.Suite
	mockDatabase *mocks.Database
	fixtureStatusRepository FixtureStatus
}

func TestFixtureStatusRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureStatusRepositoryTestSuite))
}

func (s *fixtureStatusRepositoryTestSuite) SetupTest() {
	s.mockDatabase = &mocks.Database{}
	s.fixtureStatusRepository = NewFixtureStatus(s.mockDatabase)
}

func (s *fixtureStatusRepositoryTestSuite) TestGetAllSuccess() {
	var statuses []model.FixtureStatus
	s.mockDatabase.EXPECT().GetAll(&statuses).Return(core.DatabaseResult{RowsAffected: 1})

	res, err := s.fixtureStatusRepository.GetAll()

	s.mockDatabase.AssertCalled(s.T(), "GetAll", &statuses)
	s.Equal(statuses, res)
	s.Nil(err)
}

func (s *fixtureStatusRepositoryTestSuite) TestGetAllError() {
	var statuses []model.FixtureStatus
	s.mockDatabase.EXPECT().GetAll(&statuses).Return(core.DatabaseResult{Error: errors.New("test")})

	res, err := s.fixtureStatusRepository.GetAll()

	s.mockDatabase.AssertCalled(s.T(), "GetAll", &statuses)
	s.Nil(res)
	s.ErrorContains(err, "test")
}

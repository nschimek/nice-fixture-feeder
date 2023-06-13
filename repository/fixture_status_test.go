package repository

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type fixtureStatusRepositoryTestSuite struct {
	suite.Suite
	mockDatabase *core.MockDatabase
	fixtureStatusRepository FixtureStatusRepository
}

func TestFixtureStatusRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureStatusRepositoryTestSuite))
}

func (s *fixtureStatusRepositoryTestSuite) SetupTest() {
	s.mockDatabase = &core.MockDatabase{}
	s.fixtureStatusRepository = NewFixtureStatusRepository(s.mockDatabase)
}

func (s *fixtureStatusRepositoryTestSuite) TestGetAll() {
	var statues []model.FixtureStatus
	s.mockDatabase.EXPECT().GetAll(&statues)

	s.fixtureStatusRepository.GetAll()

	s.mockDatabase.AssertCalled(s.T(), "GetAll", &statues)
}

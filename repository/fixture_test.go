package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type fixtureRepositoryTestSuite struct {
	suite.Suite
	fixtures []model.Fixture
	mockDatabase *mocks.Database
	repo *Fixture
}

func TestFixtureRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureRepositoryTestSuite))
}

func (s *fixtureRepositoryTestSuite) SetupTest() {
	s.fixtures = []model.Fixture{
		{
			Fixture: model.FixtureFixture{Id: 100, Date: time.Date(2023, 3, 5, 16, 30, 0, 0, core.UTC)},
			League: model.FixtureLeague{Id: 39, Season: 2022},
			Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 40}, Away: model.FixtureTeam{Id: 33}},
			Goals: model.FixtureGoals{Home: 7, Away: 0},
		},
	}
	s.mockDatabase = &mocks.Database{}
	s.repo = NewFixture(s.mockDatabase)
}

func (s *fixtureRepositoryTestSuite) TestUpsertSuccess() {
	r := core.DatabaseResult{RowsAffected: 1, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.fixtures).Return(r)

	actual, err := s.repo.Upsert(s.fixtures)

	s.Equal(s.fixtures, actual)
	s.Nil(err)
}

func (s *fixtureRepositoryTestSuite) TestUpsertError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.fixtures).Return(r)

	actual, err := s.repo.Upsert(s.fixtures)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

func (s *fixtureRepositoryTestSuite) TestUpsertEmptyAndNil() {
	a1, e1 := s.repo.Upsert([]model.Fixture{})
	s.Nil(a1)
	s.Nil(e1)
	a2, e2 := s.repo.Upsert(nil)
	s.Nil(a2)
	s.Nil(e2)
}
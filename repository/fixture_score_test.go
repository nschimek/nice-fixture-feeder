package repository

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

type fixtureScoreRepositoryTestSuite struct {
	suite.Suite
	fixtureScores []model.FixtureScore
	mockDatabase *mocks.Database
	repo         *FixtureScore
}

func TestFixtureScoreRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureScoreRepositoryTestSuite))
}

func (s *fixtureScoreRepositoryTestSuite) SetupTest() {
	s.mockDatabase = new(mocks.Database)
	s.fixtureScores = []model.FixtureScore{
		{FixtureId: 100, ScoreId: 1, Value: 87},
		{FixtureId: 100, ScoreId: 2, Value: 42},
		{FixtureId: 100, ScoreId: 3, Value: 12},
	}
	s.repo = NewFixtureScore(s.mockDatabase)
}

func (s *fixtureScoreRepositoryTestSuite) TestUpsertSuccess() {
	r := core.DatabaseResult{RowsAffected: 3, Error: nil}

	s.mockDatabase.EXPECT().Upsert(&s.fixtureScores).Return(r)

	actual, err := s.repo.Upsert(s.fixtureScores)

	s.Equal(s.fixtureScores, actual)
	s.Nil(err)
}

func (s *fixtureScoreRepositoryTestSuite) TestUpsertError() {
	r := core.DatabaseResult{RowsAffected: 0, Error: errors.New("test")}

	s.mockDatabase.EXPECT().Upsert(&s.fixtureScores).Return(r)

	actual, err := s.repo.Upsert(s.fixtureScores)

	s.Nil(actual)
	s.ErrorContains(err, "test")
}

func (s *fixtureScoreRepositoryTestSuite) TestUpsertEmptyAndNil() {
	a1, e1 := s.repo.Upsert([]model.FixtureScore{})
	s.Nil(a1)
	s.Nil(e1)
	a2, e2 := s.repo.Upsert(nil)
	s.Nil(a2)
	s.Nil(e2)
}
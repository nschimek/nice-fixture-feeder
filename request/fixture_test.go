package request

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/stretchr/testify/suite"
)

type fixtureRequestTestSuite struct {
	suite.Suite
	mockRequest *MockRequester[model.Fixture]
	mockRepository *repository.MockRepository[model.Fixture]
	fixtureRequest *fixtureRequest
	fixtures []model.Fixture
}

func TestFixtureRequestTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureRequestTestSuite))
}

func (s *fixtureRequestTestSuite) SetupTest() {
	s.mockRequest = &MockRequester[model.Fixture]{}
	s.mockRepository = &repository.MockRepository[model.Fixture]{}
	s.fixtureRequest = &fixtureRequest{
		config: core.MockConfig, 
		requester: s.mockRequest,
		repo: s.mockRepository,
	}
	s.fixtures = []model.Fixture{
		{
			Fixture: model.FixtureFixture{Id: 100, Date: time.Date(2023, 3, 5, 16, 30, 0, 0, core.UTC)},
			League: model.FixtureLeague{Id: 39, Season: 2022},
			Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 40}, Away: model.FixtureTeam{Id: 33}},
			Goals: model.FixtureGoals{Home: 7, Away: 0},
		},
	}
}

func (s *fixtureRequestTestSuite) TestRequestValid() {
	// expected parameters
	p := url.Values{"league": {"39"}, "season": {"2022"}, "from": {"2023-03-05"}, "to": {"2023-03-06"}}
	// response
	r := &Response[model.Fixture]{Response: s.fixtures}

	s.mockRequest.EXPECT().Get(fixturesEndpoint, p).Return(r, nil)

	s.fixtureRequest.Request(time.Date(2023, 3, 5, 0, 0, 0, 0, core.UTC), time.Date(2023, 3, 6, 0, 0, 0, 0, core.UTC), "39")

	s.Contains(s.fixtureRequest.GetData(), r.Response[0])
}

func (s *fixtureRequestTestSuite) TestRequestError() {
	p := url.Values{"league": {"39"}, "season": {"2022"}, "from": {"2023-03-05"}, "to": {"2023-03-06"}}
	s.mockRequest.EXPECT().Get(fixturesEndpoint, p).Return(nil, errors.New("test"))

	s.fixtureRequest.Request(time.Date(2023, 3, 5, 0, 0, 0, 0, core.UTC), time.Date(2023, 3, 6, 0, 0, 0, 0, core.UTC), "39")

	s.Len(s.fixtureRequest.GetData(), 0)
}

func (s *fixtureRequestTestSuite) TestPersistSuccess() {
	rs := &repository.ResultStats{
		Success: map[string]int{"fixture": 1},
		Error: map[string]int{"fixture": 0},
	}
	// prep by pre-populating with leagues, and mocking the Upsert result stats
	s.fixtureRequest.requestedData = s.fixtures
	s.mockRepository.EXPECT().Upsert(s.fixtures).Return(rs)

	s.fixtureRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", s.fixtures)
}

func (s *fixtureRequestTestSuite) TestPersistError() {
	rs := &repository.ResultStats{
		Success: map[string]int{"fixture": 0},
		Error: map[string]int{"fixture": 1},
	}
	// prep by pre-populating with leagues, and mocking the Upsert result stats
	s.fixtureRequest.requestedData = s.fixtures
	s.mockRepository.EXPECT().Upsert(s.fixtures).Return(rs)

	s.fixtureRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", s.fixtures)
	s.Nil(s.fixtureRequest.requestedData)
}

func (s *fixtureRequestTestSuite) TestPostPersist() {
	s.fixtureRequest.requestedData = s.fixtures
	s.fixtureRequest.postPersist()

	s.Contains(s.fixtureRequest.GetIds(), 100)
	s.Equal(s.fixtureRequest.GetById(100), &s.fixtures[0])
	s.Nil(s.fixtureRequest.GetById(101))
}
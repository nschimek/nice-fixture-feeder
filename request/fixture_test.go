package request

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/nschimek/nice-fixture-feeder/core"
	core_mocks "github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/nschimek/nice-fixture-feeder/request/mocks"
	"github.com/stretchr/testify/suite"
)

type fixtureRequestTestSuite struct {
	suite.Suite
	mockRequest *mocks.Requester[model.Fixture]
	mockRepository *repo_mocks.UpsertRepository[model.Fixture]
	fixtureRequest *fixture
	fixtures []model.Fixture
}

func TestFixtureRequestTestSuite(t *testing.T) {
	suite.Run(t, new(fixtureRequestTestSuite))
}

func (s *fixtureRequestTestSuite) SetupTest() {
	// use a slightly different config for these tests (make a copy so we don't impact others!)
	cfg := core_mocks.Config
	cfg.Leagues = []int{39}

	s.mockRequest = &mocks.Requester[model.Fixture]{}
	s.mockRepository = &repo_mocks.UpsertRepository[model.Fixture]{}
	s.fixtureRequest = &fixture{
		config: &cfg, 
		requester: s.mockRequest,
		repo: s.mockRepository,
		fixtureMap: make(map[int]model.Fixture),
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
	r := &model.Response[model.Fixture]{Response: s.fixtures}

	s.mockRequest.EXPECT().Get(fixturesEndpoint, p).Return(r, nil)

	s.fixtureRequest.RequestDateRange(time.Date(2023, 3, 5, 0, 0, 0, 0, core.UTC), time.Date(2023, 3, 6, 0, 0, 0, 0, core.UTC))

	s.Contains(s.fixtureRequest.requestedData, r.Response[0])
}

func (s *fixtureRequestTestSuite) TestRequestValidNoDate() {
		// expected parameters
		p := url.Values{"league": {"39"}, "season": {"2022"}}
		// response
		r := &model.Response[model.Fixture]{Response: s.fixtures}
	
		s.mockRequest.EXPECT().Get(fixturesEndpoint, p).Return(r, nil)
	
		s.fixtureRequest.Request()
		s.fixtureRequest.RequestDateRange(time.Date(2023, 3, 5, 0, 0, 0, 0, core.UTC), time.Time{})
		s.fixtureRequest.RequestDateRange(time.Time{}, time.Date(2023, 3, 5, 0, 0, 0, 0, core.UTC))
	
		s.mockRequest.AssertCalled(s.T(), "Get", fixturesEndpoint, p)
}

func (s *fixtureRequestTestSuite) TestRequestError() {
	p := url.Values{"league": {"39"}, "season": {"2022"}, "from": {"2023-03-05"}, "to": {"2023-03-06"}}
	s.mockRequest.EXPECT().Get(fixturesEndpoint, p).Return(nil, errors.New("test"))

	s.fixtureRequest.RequestDateRange(time.Date(2023, 3, 5, 0, 0, 0, 0, core.UTC), time.Date(2023, 3, 6, 0, 0, 0, 0, core.UTC))

	s.Len(s.fixtureRequest.requestedData, 0)
}

func (s *fixtureRequestTestSuite) TestPersistSuccess() {
	s.fixtureRequest.requestedData = s.fixtures
	s.mockRepository.EXPECT().Upsert(s.fixtures).Return(s.fixtures, nil)

	s.fixtureRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", s.fixtures)
	// check postPersist results
	s.Contains(s.fixtureRequest.GetIds(), 100)
	s.Contains(s.fixtureRequest.GetMap(), 100)
	s.NotContains(s.fixtureRequest.GetIds(), 101)
	s.NotContains(s.fixtureRequest.GetMap(), 101)
}

func (s *fixtureRequestTestSuite) TestPersistError() {
	s.fixtureRequest.requestedData = s.fixtures
	s.mockRepository.EXPECT().Upsert(s.fixtures).Return(nil, errors.New("test"))

	s.fixtureRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", s.fixtures)
	s.Nil(s.fixtureRequest.requestedData)
	s.Len(s.fixtureRequest.GetIds(), 0)
	s.Len(s.fixtureRequest.GetMap(), 0)
}
package request

import (
	"errors"
	"net/url"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service"
	"github.com/stretchr/testify/suite"
)

type leagueRequestTestSuite struct {
	suite.Suite
	mockRequest *MockRequester[model.League]
	mockRepository *repository.MockUpsertRepository[model.League]
	mockImageService *service.MockImageService
	leagueRequest *leagueRequest
	leagues []model.League
}

func TestLeagueRequestTestSuite(t *testing.T) {
	suite.Run(t, new(leagueRequestTestSuite))
}

func (s *leagueRequestTestSuite) SetupTest() {
	s.mockRequest = &MockRequester[model.League]{}
	s.mockRepository = &repository.MockUpsertRepository[model.League]{}
	s.mockImageService = &service.MockImageService{}
	s.leagueRequest = &leagueRequest{
		config: core.MockConfig, 
		requester: s.mockRequest,
		repo: s.mockRepository,
		imageService: s.mockImageService,
	}
	s.leagues = []model.League{
		{
			League: model.LeagueLeague{Id: 39, Name: "Premier League", Logo: "39.png"},
			Country: model.LeagueCountry{Flag: "eng.png"},
		},
		{
			League: model.LeagueLeague{Id: 140, Name: "La Liga"},
			Country: model.LeagueCountry{Flag: "spn.png"},
		},
	}
}

func (s *leagueRequestTestSuite) TestRequestValid() {
	// parameters
	p1 := url.Values{"id": {"39"}, "season": {"2022"}}
	p2 := url.Values{"id": {"140"}, "season": {"2022"}}
	// responses
	r1 := &Response[model.League]{Response: []model.League{s.leagues[0]}}
	r2 := &Response[model.League]{Response: []model.League{s.leagues[1]}}

	s.mockRequest.EXPECT().Get(leaguesEndpoint, p1).Return(r1, nil)
	s.mockRequest.EXPECT().Get(leaguesEndpoint, p2).Return(r2, nil)

	s.leagueRequest.Request()

	s.Contains(s.leagueRequest.requestedData, r1.Response[0])
	s.Contains(s.leagueRequest.requestedData, r2.Response[0])
}

func (s *leagueRequestTestSuite) TestRequestError() {
	// parameters
	p1 := url.Values{"id": {"39"}, "season": {"2022"}}
	p2 := url.Values{"id": {"140"}, "season": {"2022"}}
	s.mockRequest.EXPECT().Get(leaguesEndpoint, p1).Return(nil, errors.New("test"))
	s.mockRequest.EXPECT().Get(leaguesEndpoint, p2).Return(nil, errors.New("test"))

	s.leagueRequest.Request()

	s.Len(s.leagueRequest.requestedData, 0)
}

func (s *leagueRequestTestSuite) TestPersist() {
	s.leagueRequest.requestedData = s.leagues

	s.mockRepository.EXPECT().Upsert(s.leagues).Return(s.leagues, nil)
	// postPersist
	s.mockImageService.EXPECT().TransferURL(s.leagues[0].League.Logo, core.MockConfig.AWS.BucketName, leagueKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.leagues[0].Country.Flag, core.MockConfig.AWS.BucketName, countryKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.leagues[1].League.Logo, core.MockConfig.AWS.BucketName, leagueKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.leagues[1].Country.Flag, core.MockConfig.AWS.BucketName, countryKeyFormat).Return(true)

	s.leagueRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", s.leagues)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.leagues[0].League.Logo, core.MockConfig.AWS.BucketName, leagueKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.leagues[0].Country.Flag, core.MockConfig.AWS.BucketName, countryKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.leagues[1].League.Logo, core.MockConfig.AWS.BucketName, leagueKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.leagues[1].Country.Flag, core.MockConfig.AWS.BucketName, countryKeyFormat)
}

func (s *leagueRequestTestSuite) TestPersistError() {
	s.leagueRequest.requestedData = s.leagues

	s.mockRepository.EXPECT().Upsert(s.leagues).Return(nil, errors.New("test"))

	s.leagueRequest.Persist()

	s.Nil(s.leagueRequest.requestedData)
	s.mockImageService.AssertNotCalled(s.T(), "TransferURL")
}
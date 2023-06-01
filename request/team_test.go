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

type teamRequestTestSuite struct {
	suite.Suite
	mockRequest *MockRequester[model.Team]
	mockRepository *repository.MockRepository[model.Team]
	mockImageService *service.MockImageService
	teamRequest *teamRequest
	teams []model.Team
}

func TestTeamRequestTestSuite(t *testing.T) {
	suite.Run(t, new(teamRequestTestSuite))
}

func (s *teamRequestTestSuite) SetupTest() {
	s.mockRequest = &MockRequester[model.Team]{}
	s.mockRepository = &repository.MockRepository[model.Team]{}
	s.mockImageService = &service.MockImageService{}
	s.teamRequest = &teamRequest{
		config: core.MockConfig, 
		requester: s.mockRequest,
		repo: s.mockRepository,
		imageService: s.mockImageService,
	}
	s.teams = []model.Team{
		{Team: model.TeamTeam{Id: 40, Name: "Liverpool", Logo: "40.png"}},
		{Team: model.TeamTeam{Id: 33, Name: "Manchester United", Logo: "33.png"}},
		{Team: model.TeamTeam{Id: 42, Name: "Arsenal", Logo: "42.png"}},
		{Team: model.TeamTeam{Id: 529, Name: "Barcelona", Logo: "529.png"}},
		{Team: model.TeamTeam{Id: 530, Name: "Athletico Madrid", Logo: "530.png"}},
		{Team: model.TeamTeam{Id: 541, Name: "Real Madrid", Logo: "541.png"}},
	}
}

func (s *teamRequestTestSuite) TestRequestValid() {
		// parameters
		p1 := url.Values{"league": {"39"}, "season": {"2022"}}
		p2 := url.Values{"league": {"140"}, "season": {"2022"}}
		// responses
		r1 := &Response[model.Team]{Response: s.teams[0:3]}
		r2 := &Response[model.Team]{Response: s.teams[3:6]}

		s.mockRequest.EXPECT().Get(teamsEndpoint, p1).Return(r1, nil)
		s.mockRequest.EXPECT().Get(teamsEndpoint, p2).Return(r2, nil)

		s.teamRequest.Request(map[string]struct{}{"39": {}, "140": {}})

		s.Len(s.teamRequest.GetData(), 6)
		s.Equal(s.teamRequest.GetData()[0].Team.Name, "Liverpool")
		s.Equal(s.teamRequest.GetData()[0].TeamLeagueSeason.LeagueId, 39)
		s.Equal(s.teamRequest.GetData()[0].TeamLeagueSeason.Season, 2022)
		s.Equal(s.teamRequest.GetData()[3].Team.Name, "Barcelona")
		s.Equal(s.teamRequest.GetData()[3].TeamLeagueSeason.LeagueId, 140)
		s.Equal(s.teamRequest.GetData()[3].TeamLeagueSeason.Season, 2022)
}

func (s *teamRequestTestSuite) TestRequestError() {
	p := url.Values{"league": {"39"}, "season": {"2022"}}
	s.mockRequest.EXPECT().Get(teamsEndpoint, p).Return(nil, errors.New("test"))

	s.teamRequest.Request(map[string]struct{}{"39": {}})

	s.Len(s.teamRequest.GetData(), 0)
}

func (s *teamRequestTestSuite) TestPersist() {
	rs := &repository.ResultStats{
		Success: map[string]int{"team": 6, "team_league_season": 6},
		Error: map[string]int{"team": 0, "team_league_season": 0},
	}
	// prep by pre-populating with leagues, and mocking the Upsert result stats
	s.teamRequest.RequestedData = s.teams
	s.mockRepository.EXPECT().Upsert(s.teams).Return(rs)

	s.teamRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", s.teams)
}

func (s *teamRequestTestSuite) TestPostPersist() {
	s.teamRequest.RequestedData = s.teams[0:3]

	s.mockImageService.EXPECT().TransferURL(s.teams[0].Team.Logo, core.MockConfig.AWS.BucketName, teamKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.teams[1].Team.Logo, core.MockConfig.AWS.BucketName, teamKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.teams[2].Team.Logo, core.MockConfig.AWS.BucketName, teamKeyFormat).Return(true)

	s.teamRequest.PostPersist()

	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.teams[0].Team.Logo, core.MockConfig.AWS.BucketName, teamKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.teams[1].Team.Logo, core.MockConfig.AWS.BucketName, teamKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.teams[2].Team.Logo, core.MockConfig.AWS.BucketName, teamKeyFormat)
}
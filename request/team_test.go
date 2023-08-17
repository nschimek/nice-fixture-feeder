package request

import (
	"errors"
	"net/url"
	"testing"

	core_mocks "github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/nschimek/nice-fixture-feeder/request/mocks"
	svc_mocks "github.com/nschimek/nice-fixture-feeder/service/mocks"
	"github.com/stretchr/testify/suite"
)

type teamRequestTestSuite struct {
	suite.Suite
	mockRequest *mocks.Requester[model.Team]
	mockRepository *repo_mocks.UpsertRepository[model.Team]
	mockImageService *svc_mocks.Image
	teamRequest *team
	teams []model.Team
}

func TestTeamRequestTestSuite(t *testing.T) {
	suite.Run(t, new(teamRequestTestSuite))
}

func (s *teamRequestTestSuite) SetupTest() {
	s.mockRequest = &mocks.Requester[model.Team]{}
	s.mockRepository = &repo_mocks.UpsertRepository[model.Team]{}
	s.mockImageService = &svc_mocks.Image{}
	s.teamRequest = &team{
		config: &core_mocks.Config, 
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
		r1 := &model.Response[model.Team]{Response: s.teams[0:3]}
		r2 := &model.Response[model.Team]{Response: s.teams[3:6]}

		s.mockRequest.EXPECT().Get(teamsEndpoint, p1).Return(r1, nil)
		s.mockRequest.EXPECT().Get(teamsEndpoint, p2).Return(r2, nil)

		s.teamRequest.Request()

		s.Len(s.teamRequest.requestedData, 6)
		s.Contains(s.teamRequest.requestedData, model.Team{
			Team: s.teams[0].Team, 
			TeamLeagueSeason: model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}},
		})
		s.Contains(s.teamRequest.requestedData, model.Team{
			Team: s.teams[3].Team, 
			TeamLeagueSeason: model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 529, LeagueId: 140, Season: 2022}},
		})
}

func (s *teamRequestTestSuite) TestRequestError() {
	p1 := url.Values{"league": {"39"}, "season": {"2022"}}
	p2 := url.Values{"league": {"140"}, "season": {"2022"}}
	s.mockRequest.EXPECT().Get(teamsEndpoint, p1).Return(nil, errors.New("test"))
	s.mockRequest.EXPECT().Get(teamsEndpoint, p2).Return(nil, errors.New("test"))

	s.teamRequest.Request()

	s.Len(s.teamRequest.requestedData, 0)
}

func (s *teamRequestTestSuite) TestPersistSuccess() {
	teams := s.teams[0:3] // don't need all of them for this test
	s.teamRequest.requestedData = teams

	s.mockRepository.EXPECT().Upsert(teams).Return(teams, nil)
	// postPersist
	s.mockImageService.EXPECT().TransferURL(s.teams[0].Team.Logo, core_mocks.Config.AWS.BucketName, teamKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.teams[1].Team.Logo, core_mocks.Config.AWS.BucketName, teamKeyFormat).Return(true)
	s.mockImageService.EXPECT().TransferURL(s.teams[2].Team.Logo, core_mocks.Config.AWS.BucketName, teamKeyFormat).Return(true)

	s.teamRequest.Persist()

	s.mockRepository.AssertCalled(s.T(), "Upsert", teams)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.teams[0].Team.Logo, core_mocks.Config.AWS.BucketName, teamKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.teams[1].Team.Logo, core_mocks.Config.AWS.BucketName, teamKeyFormat)
	s.mockImageService.AssertCalled(s.T(), "TransferURL", s.teams[2].Team.Logo, core_mocks.Config.AWS.BucketName, teamKeyFormat)
}

func (s *teamRequestTestSuite) TestPersistError() {
	s.teamRequest.requestedData = s.teams

	s.mockRepository.EXPECT().Upsert(s.teams).Return(nil, errors.New("test"))

	s.teamRequest.Persist()

	s.Nil(s.teamRequest.requestedData)
	s.mockImageService.AssertNotCalled(s.T(), "TransferURL")
}
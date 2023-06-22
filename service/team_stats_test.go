package service

import (
	"testing"

	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/stretchr/testify/suite"
)

type teamStatsServiceTestSuite struct {
	suite.Suite
	mockTeamStatsRepo *repository.MockTeamStatsRepository
	mockTlsRepo *repository.MockTeamLeagueSeasonRepository
	mockStatusService *MockFixtureStatusService
	teamStatsService *teamStatsService
	fixtures []model.Fixture
	fixtureIds []int
}

func TestTeamStatsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(teamStatsServiceTestSuite))
}

func (s *teamStatsServiceTestSuite) SetupTest() {
	s.mockTeamStatsRepo = &repository.MockTeamStatsRepository{}
	s.mockTlsRepo = &repository.MockTeamLeagueSeasonRepository{}
	s.mockStatusService = &MockFixtureStatusService{}
	s.teamStatsService = NewTeamStatsService(s.mockTeamStatsRepo, s.mockTlsRepo, s.mockStatusService)
	s.fixtures = []model.Fixture{
		{
			Fixture: model.FixtureFixture{Id: 100, Status: model.FixtureStatusId{Id: "FT"}},
			League: model.FixtureLeague{Id: 39, Season: 2022},
			Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 31, Result: "L"}, Away: model.FixtureTeam{Id: 40, Result: "W"}},
			Goals: model.FixtureGoals{Home: 1, Away: 6},
		},
		{
			Fixture: model.FixtureFixture{Id: 101, Status: model.FixtureStatusId{Id: "FT"}},
			League: model.FixtureLeague{Id: 39, Season: 2022},
			Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 40, Result: "L"}, Away: model.FixtureTeam{Id: 27, Result: "W"}},
			Goals: model.FixtureGoals{Home: 2, Away: 3},
		},
		{
			Fixture: model.FixtureFixture{Id: 102, Status: model.FixtureStatusId{Id: "FT"}},
			League: model.FixtureLeague{Id: 39, Season: 2022},
			Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 40, Result: "D"}, Away: model.FixtureTeam{Id: 41, Result: "D"}},
			Goals: model.FixtureGoals{Home: 0, Away: 0},
		},
		{
			Fixture: model.FixtureFixture{Id: 103, Status: model.FixtureStatusId{Id: "NS"}},
			League: model.FixtureLeague{Id: 39, Season: 2022},
			Teams: model.FixtureTeams{Home: model.FixtureTeam{Id: 38}, Away: model.FixtureTeam{Id: 40}},
		},
	}
	s.fixtureIds = []int{100, 101, 102, 103}
}

func (s *teamStatsServiceTestSuite) TestMaintainStats() {
	s.mockStatusService.EXPECT().IsFinished("FT").Return(true)
}

func (s *teamStatsServiceTestSuite) TestCalculateCurrentStats() {
	prev := &model.TeamStats{TeamStatsId: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 99}}
	a1, err1 := s.teamStatsService.calculateCurrentStats(prev, &s.fixtures[0])

	// expected (1st iteration)
	e1 := &model.TeamStats{
		TeamStatsId: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100},
		TeamStatsFixtures: model.TeamStatsFixtures{
			FixturesPlayed: model.TeamStatsHomeAwayTotal{Home: 0, Away: 1, Total: 1},
			FixturesWins: model.TeamStatsHomeAwayTotal{Home: 0, Away: 1, Total: 1},
			FixturesLosses: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
			FixturesDraws: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
		},
		TeamStatsGoals: model.TeamStatsGoals{
			GoalsFor: model.TeamStatsHomeAwayTotal{Home: 0, Away: 6, Total: 6},
			GoalsAgainst: model.TeamStatsHomeAwayTotal{Home: 0, Away: 1, Total: 1},
		},
		Form: "W",
		GoalDifferential: 5,
		CleanSheets: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
		FailedToScore: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
	}

	s.Equal(e1, a1)
	s.Nil(err1)

	a2, err2 := s.teamStatsService.calculateCurrentStats(a1, &s.fixtures[1])

	// expected (2nd iteration)
	e2 := &model.TeamStats{
		TeamStatsId: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 101},
		TeamStatsFixtures: model.TeamStatsFixtures{
			FixturesPlayed: model.TeamStatsHomeAwayTotal{Home: 1, Away: 1, Total: 2},
			FixturesWins: model.TeamStatsHomeAwayTotal{Home: 0, Away: 1, Total: 1},
			FixturesLosses: model.TeamStatsHomeAwayTotal{Home: 1, Away: 0, Total: 1},
			FixturesDraws: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
		},
		TeamStatsGoals: model.TeamStatsGoals{
			GoalsFor: model.TeamStatsHomeAwayTotal{Home: 2, Away: 6, Total: 8},
			GoalsAgainst: model.TeamStatsHomeAwayTotal{Home: 3, Away: 1, Total: 4},
		},
		Form: "WL",
		GoalDifferential: 4,
		CleanSheets: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
		FailedToScore: model.TeamStatsHomeAwayTotal{Home: 0, Away: 0, Total: 0},
	}

	s.Equal(e2, a2)
	s.Nil(err2)

	a3, err3 := s.teamStatsService.calculateCurrentStats(a2, &s.fixtures[2])

	e3 := &model.TeamStats{
		TeamStatsId: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 102},
		TeamStatsFixtures: model.TeamStatsFixtures{
			FixturesPlayed: model.TeamStatsHomeAwayTotal{Home: 2, Away: 1, Total: 3},
			FixturesWins: model.TeamStatsHomeAwayTotal{Home: 0, Away: 1, Total: 1},
			FixturesLosses: model.TeamStatsHomeAwayTotal{Home: 1, Away: 0, Total: 1},
			FixturesDraws: model.TeamStatsHomeAwayTotal{Home: 1, Away: 0, Total: 1},
		},
		TeamStatsGoals: model.TeamStatsGoals{
			GoalsFor: model.TeamStatsHomeAwayTotal{Home: 2, Away: 6, Total: 8},
			GoalsAgainst: model.TeamStatsHomeAwayTotal{Home: 3, Away: 1, Total: 4},
		},
		Form: "WLD",
		GoalDifferential: 4,
		CleanSheets: model.TeamStatsHomeAwayTotal{Home: 1, Away: 0, Total: 1},
		FailedToScore: model.TeamStatsHomeAwayTotal{Home: 1, Away: 0, Total: 1},
	}

	s.Equal(e3, a3)
	s.Nil(err3)
}

func (s *teamStatsServiceTestSuite) TestCalculateCurrentStatsError() {
	prev := &model.TeamStats{TeamStatsId: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 101}}
	_, err := s.teamStatsService.calculateCurrentStats(prev, &s.fixtures[0])

	s.ErrorContains(err, "GTE")
}


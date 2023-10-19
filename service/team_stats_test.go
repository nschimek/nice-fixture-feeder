package service

import (
	"errors"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/nschimek/nice-fixture-feeder/service/mocks"

	"github.com/stretchr/testify/suite"
)

type teamStatsServiceTestSuite struct {
	suite.Suite
	mockTsRepo *repo_mocks.TeamStats
	mockTlsService *mocks.TeamLeagueSeason
	mockStatusService *mocks.FixtureStatus
	teamStatsService *teamStats
	fixtures []model.Fixture
	fixtureIds []int
}

func TestTeamStatsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(teamStatsServiceTestSuite))
}

func (s *teamStatsServiceTestSuite) SetupTest() {
	s.mockTsRepo = &repo_mocks.TeamStats{}
	s.mockTlsService = &mocks.TeamLeagueSeason{}
	s.mockStatusService = &mocks.FixtureStatus{}
	s.teamStatsService = NewTeamStats(s.mockTsRepo, s.mockTlsService, s.mockStatusService)
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

func (s *teamStatsServiceTestSuite) TestGetById() {
	ts := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100, NextFixtureId: 101}}

	s.mockTsRepo.EXPECT().GetById(ts.Id).Return(&ts, nil)

	res, err := s.teamStatsService.GetById(ts.Id)

	s.Nil(err)
	s.Equal(&ts, res)
	s.Contains(s.teamStatsService.statsMap, ts.Id.GetCurrentId())
	s.Contains(s.teamStatsService.statsMap, ts.Id.GetNextId())
}

func (s *teamStatsServiceTestSuite) TestGetByIdFromMap() {
	ts := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100, NextFixtureId: 101}, Form: "W"}
	s.teamStatsService.statsMap[ts.Id.GetCurrentId()] = ts

	res, err := s.teamStatsService.GetById(ts.Id.GetCurrentId())

	s.Nil(err)
	s.Equal(&ts, res)
	s.Contains(s.teamStatsService.statsMap, ts.Id.GetCurrentId())
	s.Contains(s.teamStatsService.statsMap, ts.Id.GetNextId())
}

func (s *teamStatsServiceTestSuite) TestMaintainStats() {
	tlsHome := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 31, LeagueId: 39, Season: 2022}}
	tlsAway := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}

	s.mockStatusService.EXPECT().IsFinished("FT").Return(true)
	s.mockStatusService.EXPECT().IsFinished("NS").Return(false)
	s.mockTlsService.EXPECT().GetById(tlsHome).Return(&tlsHome, nil)
	s.mockTlsService.EXPECT().GetById(tlsAway).Return(&tlsAway, nil)
	s.mockTsRepo.AssertNotCalled(s.T(), "GetById") // TLS has max fixture ID of 0, so this should not be called

	// test with one completed fixture, one not started, and one ID not in the map (to cover all branches)
	s.teamStatsService.MaintainStats([]int{s.fixtureIds[0], s.fixtureIds[1], s.fixtureIds[3]}, 
		map[int]model.Fixture{100: s.fixtures[0], 103: s.fixtures[3]})

	// assert that the invalid ID and the not started fixtures are not present in the results (with the len checks)
	s.Len(s.teamStatsService.tlsMap, 2)
	s.Equal(100, s.teamStatsService.tlsMap[tlsHome.Id].MaxFixtureId)
	s.Equal(100, s.teamStatsService.tlsMap[tlsAway.Id].MaxFixtureId)
	s.Len(s.teamStatsService.statsMap, 2)
	s.Contains(s.teamStatsService.statsMap, model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100})
	s.Contains(s.teamStatsService.statsMap, model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100})
}

func (s *teamStatsServiceTestSuite) TestMaintainFixtureWithPrevious() {
	f := &s.fixtures[1]
	tlsCurr := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tlsPrev := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100}
	tsid := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}

	s.mockTlsService.EXPECT().GetById(tlsCurr).Return(&tlsPrev, nil)
	s.mockTsRepo.EXPECT().GetById(model.TeamStats{Id: tsid}).Return(&model.TeamStats{Id: tsid}, nil)

	s.teamStatsService.maintainFixture(f, true)

	s.Len(s.teamStatsService.tlsMap, 1)
	s.Equal(101, s.teamStatsService.tlsMap[tlsCurr.Id].MaxFixtureId)
	// this will have to be re-worked?
	// s.Equal(101, s.teamStatsService.statsMap[model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}].Id.NextFixtureId)
	s.Contains(s.teamStatsService.statsMap, model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100})
	s.Contains(s.teamStatsService.statsMap, model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 101})
}

// Test a common possiblity: the curent MaxFixtureId is GTE the incoming one (this can happen on re-runs)
func (s *teamStatsServiceTestSuite) TestMaintainFixturePrevIdHigher() {
	f := &s.fixtures[1]
	tlsCurr := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tlsPrev := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 999}

	s.mockTlsService.EXPECT().GetById(tlsCurr).Return(&tlsPrev, nil)

	s.teamStatsService.maintainFixture(f, true)

	// maps should not be populated because the end result should be that we skip persisting this fixture entirely
	s.Len(s.teamStatsService.tlsMap, 0)
	s.Len(s.teamStatsService.statsMap, 0)
}

func (s *teamStatsServiceTestSuite) TestMaintainFixtureErrorNoTLS() {
	f := &s.fixtures[1]
	tls := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}

	s.mockTlsService.EXPECT().GetById(tls).Return(nil, nil)

	s.teamStatsService.maintainFixture(f, true)

	// maps should not be populated due to error from getTLS().
	s.Len(s.teamStatsService.tlsMap, 0)
	s.Len(s.teamStatsService.statsMap, 0)
}

func (s *teamStatsServiceTestSuite) TestGetUpdatedStatsErrorPrevious() {
	f := &s.fixtures[0]
	tsidCurr := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}
	tlsReq := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tsidPrev := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 99}
	tlsRes := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 99}

	s.mockTlsService.EXPECT().GetById(tlsReq).Return(&tlsRes, nil)
	s.mockTsRepo.EXPECT().GetById(model.TeamStats{Id: tsidPrev}).Return(nil, nil)

	tls, curr, prev, err := s.teamStatsService.getUpdatedStats(&tsidCurr, f)

	s.Nil(tls)
	s.Nil(curr)
	s.Nil(prev)
	s.ErrorContains(err, "no stats")
}

func (s *teamStatsServiceTestSuite) TestGetUpdatedStatsErrorCurrent() {
	f := &s.fixtures[0]
	tsidCurr := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}
	tlsReq := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tsidPrev := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 99}
	tlsRes := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 99}

	s.mockTlsService.EXPECT().GetById(tlsReq).Return(&tlsRes, nil)
	s.mockTsRepo.EXPECT().GetById(model.TeamStats{Id: tsidPrev}).Return(&model.TeamStats{
		Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 999}, 
	}, nil)

	tls, curr, prev, err := s.teamStatsService.getUpdatedStats(&tsidCurr, f)

	s.Nil(tls)
	s.Nil(curr)
	s.Nil(prev)
	s.ErrorContains(err, "previous fixture ID")
}

func (s *teamStatsServiceTestSuite) TestGetPreviousStatsExisting() {
	tls := &model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100}
	tsid := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}
	ts := model.TeamStats{Id: tsid, Form: "W"}

	s.teamStatsService.statsMap[tsid] = ts
	a, err := s.teamStatsService.getPreviousStats(tls)

	s.Equal(&ts, a)
	s.Nil(err)
}

func (s *teamStatsServiceTestSuite) TestCalculateCurrentStats() {
	prev := &model.TeamStats{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 99}}
	a1, err1 := s.teamStatsService.calculateCurrentStats(prev, &s.fixtures[0])

	// expected (1st iteration)
	e1 := &model.TeamStats{
		Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100},
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
		Points: 3,
	}

	s.Equal(e1, a1)
	s.Nil(err1)

	a2, err2 := s.teamStatsService.calculateCurrentStats(a1, &s.fixtures[1])

	// expected (2nd iteration)
	e2 := &model.TeamStats{
		Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 101},
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
		Points: 3,
	}

	s.Equal(e2, a2)
	s.Nil(err2)

	a3, err3 := s.teamStatsService.calculateCurrentStats(a2, &s.fixtures[2])

	e3 := &model.TeamStats{
		Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 102},
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
		Points: 4,
	}

	s.Equal(e3, a3)
	s.Nil(err3)
}

func (s *teamStatsServiceTestSuite) TestCalculateCurrentStatsError() {
	prev := &model.TeamStats{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 101}}
	_, err := s.teamStatsService.calculateCurrentStats(prev, &s.fixtures[0])

	s.ErrorContains(err, "GTE")
}

func (s *teamStatsServiceTestSuite) TestPersistSuccess() {
	s.teamStatsService.statsMap = map[model.TeamStatsId]model.TeamStats{
		{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}: {Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WWW"},
		{TeamId: 41, LeagueId: 39, Season: 2022, FixtureId: 100}: {Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WLD"},
	}
	s.teamStatsService.tlsMap = map[model.TeamLeagueSeasonId]model.TeamLeagueSeason{
		{TeamId: 40, LeagueId: 39, Season: 2022}: {Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100},
		{TeamId: 41, LeagueId: 39, Season: 2022}: {Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100},
	}

	stats := []model.TeamStats{
		{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WWW"},
		{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WLD"},
	}

	// tls := []model.TeamLeagueSeason{
	// 	{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100},
	// 	{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100},
	// }

	s.mockTsRepo.EXPECT().Upsert(stats).Return(stats, nil)

	s.teamStatsService.Persist()

	s.mockTsRepo.AssertCalled(s.T(), "Upsert", stats)
	s.mockTlsService.AssertCalled(s.T(), "Persist")
}

func (s *teamStatsServiceTestSuite) TestPersistError() {
	s.teamStatsService.statsMap = map[model.TeamStatsId]model.TeamStats{
		{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}: {Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WWW"},
		{TeamId: 41, LeagueId: 39, Season: 2022, FixtureId: 100}: {Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WLD"},
	}
	stats := []model.TeamStats{
		{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WWW"},
		{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022}, Form: "WLD"},
	}

	s.mockTsRepo.EXPECT().Upsert(stats).Return(nil, errors.New("test"))

	s.teamStatsService.Persist()
	s.mockTsRepo.AssertCalled(s.T(), "Upsert", stats)
	s.mockTlsService.AssertNotCalled(s.T(), "Persist")
}
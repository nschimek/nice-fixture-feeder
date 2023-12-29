package service

import (
	"errors"
	"testing"

	core_mocks "github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	repo_mocks "github.com/nschimek/nice-fixture-feeder/repository/mocks"
	"github.com/nschimek/nice-fixture-feeder/service/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type teamStatsServiceTestSuite struct {
	suite.Suite
	mockTsRepo *repo_mocks.TeamStats
	mockCache *core_mocks.Cache[model.TeamStats]
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
	s.mockCache = &core_mocks.Cache[model.TeamStats]{}
	s.mockTlsService = &mocks.TeamLeagueSeason{}
	s.mockStatusService = &mocks.FixtureStatus{}
	s.teamStatsService = NewTeamStats(s.mockTsRepo, s.mockCache, s.mockTlsService, s.mockStatusService)
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

func (s *teamStatsServiceTestSuite) TestGetByIdCurrentCacheMiss() {
	ts := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100, NextFixtureId: 101}}

	s.mockCache.EXPECT().Get(ts.Id.GetCurrentId()).Return(nil, nil)
	s.mockTsRepo.EXPECT().GetById(ts.Id.GetCurrentId()).Return(&ts, nil)
	s.mockCache.EXPECT().Set(ts.Id.GetCurrentId(), &ts).Return(nil)

	res, err := s.teamStatsService.GetById(ts.Id, true)

	s.Nil(err)
	s.Equal(&ts, res)
	s.mockCache.AssertExpectations(s.T())
}

func (s *teamStatsServiceTestSuite) TestGetByIdFalseCacheHit() {
	ts := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100, NextFixtureId: 101}, Form: "W"}

	s.mockCache.EXPECT().Get(ts.Id.GetNextId()).Return(&ts, nil)
	s.mockTsRepo.AssertNotCalled(s.T(), "GetById", ts.Id.GetNextId())

	res, err := s.teamStatsService.GetById(ts.Id, false)

	s.Nil(err)
	s.Equal(&ts, res)
}

func (s *teamStatsServiceTestSuite) TestGetByIdNotFound() {
	id := model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100}

	s.mockCache.EXPECT().Get(id).Return(nil, nil)
	s.mockTsRepo.EXPECT().GetById(id).Return(nil, errors.New("not found"))

	res, err := s.teamStatsService.GetById(id, true)

	s.Nil(res)
	s.ErrorContains(err, "no stats")
}

func (s *teamStatsServiceTestSuite) TestGetByIdInvalid() {
	tsc := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 100}}
	tsn := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, NextFixtureId: 100}}

	res1, err1 := s.teamStatsService.GetById(tsn.Id, true)

	s.Nil(res1)
	s.ErrorContains(err1, "current is true but FixtureId is 0")

	res2, err2 := s.teamStatsService.GetById(tsc.Id, false)

	s.Nil(res2)
	s.ErrorContains(err2, "current is false but NextFixtureId is 0")
}

func (s *teamStatsServiceTestSuite) TestGetByIdWithTLSCurrent() {
	ts := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 102, NextFixtureId: 103}}
	tlsid := ts.Id.GetTlsId()
	tls := model.TeamLeagueSeason{Id: tlsid, MaxFixtureId: 102}

	s.mockTlsService.EXPECT().GetById(tlsid).Return(&tls, nil)
	s.mockCache.EXPECT().Get(ts.Id.GetCurrentId()).Return(&ts, nil)

	res, err := s.teamStatsService.GetByIdWithTLS(ts.Id, true)

	s.Nil(err)
	s.Equal(&ts, res)
}

func (s *teamStatsServiceTestSuite) TestGetByIdWithTLSNext() {
	ts := model.TeamStats{Id: model.TeamStatsId{TeamId: 31, LeagueId: 39, Season: 2022, FixtureId: 102, NextFixtureId: 103}}
	tlsid := ts.Id.GetTlsId()
	tls := model.TeamLeagueSeason{Id: tlsid, MaxFixtureId: 103}

	s.mockTlsService.EXPECT().GetById(tlsid).Return(&tls, nil)
	s.mockCache.EXPECT().Get(ts.Id.GetNextId()).Return(&ts, nil)

	res, err := s.teamStatsService.GetByIdWithTLS(ts.Id, false)

	s.Nil(err)
	s.Equal(&ts, res)
}

func (s *teamStatsServiceTestSuite) TestMaintainStats() {
	tlsHome := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 31, LeagueId: 39, Season: 2022}}
	tlsAway := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}

	s.mockStatusService.EXPECT().IsFinished("FT").Return(true)
	s.mockStatusService.EXPECT().IsFinished("NS").Return(false)
	s.mockTlsService.EXPECT().GetById(tlsHome.Id).Return(&tlsHome, nil)
	s.mockTlsService.EXPECT().GetById(tlsAway.Id).Return(&tlsAway, nil)
	s.mockTlsService.AssertNotCalled(s.T(), "GetById") // TLS has max fixture ID of 0, so this should not be called
	s.mockTlsService.EXPECT().PersistOne(&model.TeamLeagueSeason{Id: tlsHome.Id, MaxFixtureId: 100})
	s.mockTlsService.EXPECT().PersistOne(&model.TeamLeagueSeason{Id: tlsAway.Id, MaxFixtureId: 100})
	// we are not concerned with what the persist methods are being called with in this test
	s.mockCache.EXPECT().Set(mock.AnythingOfType("model.TeamStatsId"), mock.AnythingOfType("*model.TeamStats")).Return(nil)
	s.mockTsRepo.EXPECT().UpsertOne(mock.AnythingOfType("model.TeamStats")).Return(model.TeamStats{}, nil)

	// test with one completed fixture, one not started, and one ID not in the map (to cover all branches)
	s.teamStatsService.MaintainStats([]int{s.fixtureIds[0], s.fixtureIds[1], s.fixtureIds[3]}, 
		map[int]model.Fixture{100: s.fixtures[0], 103: s.fixtures[3]})

	s.mockTlsService.AssertExpectations(s.T())
}

func (s *teamStatsServiceTestSuite) TestMaintainFixtureWithPrevious() {
	f := &s.fixtures[1]
	tlsCurr := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tlsPrev := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100}
	tsid := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}

	s.mockTlsService.EXPECT().GetById(tlsCurr.Id).Return(&tlsPrev, nil)
	s.mockCache.EXPECT().Get(tsid).Return(&model.TeamStats{Id: tsid}, nil)
	s.mockTlsService.EXPECT().PersistOne(&model.TeamLeagueSeason{Id: tlsCurr.Id, MaxFixtureId: 101})

	// once again we are not concerned with the stat values here, just that persist is getting called	
	s.mockCache.EXPECT().Set(model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}, 
		mock.AnythingOfType("*model.TeamStats")).Return(nil)
	s.mockCache.EXPECT().Set(model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, NextFixtureId: 101}, 
		mock.AnythingOfType("*model.TeamStats")).Return(nil)
	s.mockCache.EXPECT().Set(model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 101}, 
		mock.AnythingOfType("*model.TeamStats")).Return(nil)
	s.mockTsRepo.EXPECT().UpsertOne(mock.AnythingOfType("model.TeamStats")).Return(model.TeamStats{}, nil)

	s.teamStatsService.maintainFixture(f, true)

	// should end up with just 2 persist calls
	s.mockTsRepo.AssertNumberOfCalls(s.T(), "UpsertOne", 2)
	s.mockCache.AssertExpectations(s.T())
}

// Test a common possiblity: the curent MaxFixtureId is GTE the incoming one (this can happen on re-runs)
func (s *teamStatsServiceTestSuite) TestMaintainFixturePrevIdHigher() {
	f := &s.fixtures[1]
	tlsCurr := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tlsPrev := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 999}

	s.mockTlsService.EXPECT().GetById(tlsCurr.Id).Return(&tlsPrev, nil)

	s.teamStatsService.maintainFixture(f, true)

	// persists should not be called
	s.mockTsRepo.AssertNotCalled(s.T(), "UpsertOne")
	s.mockTlsService.AssertNotCalled(s.T(), "PersistOne")
}

func (s *teamStatsServiceTestSuite) TestMaintainFixtureErrorNoTLS() {
	f := &s.fixtures[1]
	tlsid := model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}

	s.mockTlsService.EXPECT().GetById(tlsid).Return(nil, errors.New("no TLS"))

	s.teamStatsService.maintainFixture(f, true)

	// persists should not be called
	s.mockTsRepo.AssertNotCalled(s.T(), "UpsertOne")
	s.mockTlsService.AssertNotCalled(s.T(), "PersistOne")
}

func (s *teamStatsServiceTestSuite) TestMaintainFixtureErrorCalcCurrent() {
	f := &s.fixtures[1]
	tlsCurr := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}}
	tlsPrev := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 100}
	tsid := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}

	s.mockTlsService.EXPECT().GetById(tlsCurr.Id).Return(&tlsPrev, nil)
	s.mockCache.EXPECT().Get(tsid).Return(&model.TeamStats{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 999}}, nil)

	s.teamStatsService.maintainFixture(f, true)

	// persists should not be called due to error with calculateCurrentStats() - prev fixture ID is higher
	s.mockTsRepo.AssertNotCalled(s.T(), "UpsertOne")
	s.mockTlsService.AssertNotCalled(s.T(), "PersistOne")
}

func (s *teamStatsServiceTestSuite) TestGetUpdatedStatsErrorPrevious() {
	f := &s.fixtures[0]
	tsidCurr := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}
	tlsId := model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}
	tsidPrev := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 99}
	tlsRes := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 99}

	s.mockTlsService.EXPECT().GetById(tlsId).Return(&tlsRes, nil)
	s.mockCache.EXPECT().Get(tsidPrev).Return(nil, nil)
	s.mockTsRepo.EXPECT().GetById(tsidPrev).Return(nil, nil)

	tls, curr, prev, err := s.teamStatsService.getUpdatedStats(&tsidCurr, f)

	s.Nil(tls)
	s.Nil(curr)
	s.Nil(prev)
	s.ErrorContains(err, "no stats")
}

func (s *teamStatsServiceTestSuite) TestGetUpdatedStatsErrorCurrent() {
	f := &s.fixtures[0]
	tsidCurr := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}
	tlsId := model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}
	tsidPrev := model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 99}
	tlsRes := model.TeamLeagueSeason{Id: model.TeamLeagueSeasonId{TeamId: 40, LeagueId: 39, Season: 2022}, MaxFixtureId: 99}

	s.mockTlsService.EXPECT().GetById(tlsId).Return(&tlsRes, nil)
	s.mockCache.EXPECT().Get(tsidPrev).Return(&model.TeamStats{
		Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 999}, 
	}, nil)

	tls, curr, prev, err := s.teamStatsService.getUpdatedStats(&tsidCurr, f)

	s.Nil(tls)
	s.Nil(curr)
	s.NotNil(prev)
	s.ErrorContains(err, "previous fixture ID")
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

func (s *teamStatsServiceTestSuite) TestPersistOneBothIds() {
	ts := &model.TeamStats{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100, NextFixtureId: 101}, Form: "WWW"}

	s.mockCache.EXPECT().Set(ts.Id.GetCurrentId(), ts).Return(nil)
	s.mockCache.EXPECT().Set(ts.Id.GetNextId(), ts).Return(nil)
	s.mockTsRepo.EXPECT().UpsertOne(*ts).Return(*ts, nil)

	s.teamStatsService.PersistOne(ts)
}

func (s *teamStatsServiceTestSuite) TestPersistOneOneId() {
	ts1 := &model.TeamStats{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, FixtureId: 100}, Form: "WWW"}

	s.mockCache.EXPECT().Set(ts1.Id.GetCurrentId(), ts1).Return(nil)
	s.mockTsRepo.EXPECT().UpsertOne(*ts1).Return(*ts1, nil)

	s.teamStatsService.PersistOne(ts1)

	ts2 := &model.TeamStats{Id: model.TeamStatsId{TeamId: 40, LeagueId: 39, Season: 2022, NextFixtureId: 101}, Form: "WWW"}

	s.mockCache.EXPECT().Set(ts2.Id.GetNextId(), ts2).Return(nil)
	s.mockTsRepo.EXPECT().UpsertOne(*ts2).Return(*ts2, nil)

	s.teamStatsService.PersistOne(ts2)
}
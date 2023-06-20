package service

import (
	"errors"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/sirupsen/logrus"
)

type TeamStatsService interface {
	MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture)
}

type teamStatsService struct {
	tsRepo repository.TeamStatsRepository
	tlsRepo repository.TeamLeagueSeasonRepository
	statusService FixtureStatusService
}

func NewTeamStatsService(tsRepo repository.TeamStatsRepository, 
	tlsRepo repository.TeamLeagueSeasonRepository,
	statusService FixtureStatusService) *teamStatsService {
	return &teamStatsService{tsRepo: tsRepo, tlsRepo: tlsRepo, statusService: statusService}
}

func (s *teamStatsService) MaintainStats(fixtureIds []int, fixtureMap map[int]model.Fixture) {

}

func (s *teamStatsService) getTLS(tsid model.TeamStatsId) (*model.TeamLeagueSeason, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": tsid.TeamId, "leagueId": tsid.LeagueId, "season": tsid.Season,
	}).Debug("getting team league season (TLS)...")

	tls := s.tlsRepo.GetById(model.TeamLeagueSeason{TeamId: tsid.TeamId, LeagueId: tsid.LeagueId, Season: tsid.Season})

	if tls == nil {
		return nil, errors.New("could not get TLS, was the league setup?")
	} else {
		return tls, nil
	}
}
func (s *teamStatsService) getCurrentStats(tls model.TeamLeagueSeason, tsid model.TeamStatsId) (*model.TeamStats, error) {
	core.Log.WithFields(logrus.Fields{
		"teamId": tsid.TeamId, "leagueId": tsid.LeagueId, "season": tsid.Season,
	}).Debug("getting current stats...")
	
	stats := new(model.TeamStats)
	if tls.MaxFixtureId == 0 {
		core.Log.Debug("max fixture ID was 0, starting new stats")
		stats = &model.TeamStats{TeamStatsId: tsid}
	} else {
		stats = s.tsRepo.GetById(model.TeamStats{TeamStatsId: tsid})
		if stats == nil {
			return nil, errors.New("no stats for current max fixture ID, indicates corruption")
		}
	}

	return stats, nil
}

func (s *teamStatsService) calculateNewStats(curr model.TeamStats, fixture model.Fixture) (*model.TeamStats, error) {
	copy := curr // create a copy of the current
	upd := &copy // work with the pointer

	if (curr.TeamStatsId.Season != fixture.League.Season) {
		return nil, errors.New("current stats have different season than fixture")
	}
	if (curr.TeamStatsId.LeagueId != fixture.League.Id) {
		return nil, errors.New("current stats have different league ID than fixture")
	}

	rs := fixture.GetResultStats(curr.TeamStatsId.TeamId)

	if (rs.Result == model.ResultWin) {
		upd.TeamStatsFixtures.FixturesWins.Increment(1, rs.Home)
	} else if (rs.Result == model.ResultLoss) {
		upd.TeamStatsFixtures.FixturesLosses.Increment(1, rs.Home)
	} else {
		upd.TeamStatsFixtures.FixturesDraws.Increment(1, rs.Home)
	}
	upd.TeamStatsGoals.GoalsFor.Increment(rs.GoalsFor, rs.Home)
	upd.TeamStatsGoals.GoalsAgainst.Increment(rs.GoalsAgainst, rs.Home)

	if (rs.GoalsAgainst == 0) {
		upd.CleanSheets.Increment(1, rs.Home)
	}
	
	if (rs.GoalsFor == 0) {
		upd.FailedToScore.Increment(1, rs.Home)
	}

	upd.GoalDifferential = upd.GoalDifferential + (rs.GoalsFor - rs.GoalsAgainst)
	upd.Form = upd.Form + rs.Result

	return upd, nil
}
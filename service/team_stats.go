package service

import (
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
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

func (s *teamStatsService) getCurrentStats(tsid model.TeamStatsId) *model.TeamStats {
	tls := s.tlsRepo.GetById(model.TeamLeagueSeason{TeamId: tsid.TeamId, LeagueId: tsid.LeagueId, Season: tsid.Season})
	
	stats := new(model.TeamStats)
	if tls.MaxFixtureId == 0 {
		stats = &model.TeamStats{TeamStatsId: tsid}
	} else {
		stats = s.tsRepo.GetById(model.TeamStats{TeamStatsId: tsid})
	}

	return stats
}
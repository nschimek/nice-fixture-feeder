package service

import (
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
)

type score struct {
	tsRepo repository.TeamStats
	scores *scores.ScoreRegistry
	statusService FixtureStatus
	tlsService TeamLeagueSeason
	statsMap map[model.TeamStatsId]model.TeamStats
}

func NewScore(tsRepo repository.TeamStats, 
	scores *scores.ScoreRegistry,
	statusService FixtureStatus,
	tlsService TeamLeagueSeason) *score {
	s := &score{
		tsRepo: tsRepo, 
		scores: scores,
		statusService: statusService,
		tlsService: tlsService,
		statsMap: make(map[model.TeamStatsId]model.TeamStats),
	}
	s.setup()
	return s
}

func (s *score) setup() {
	s.scores.PointsStrength.SetStatsFunc(s.getStats)
}

func (s *score) getStats(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats) {
	// this code needs to be moved to its own function as its only for HOME
	tsid := fixture.GetTeamStatsId(true)

	if s.statusService.IsScheduled(fixture.Fixture.Status.Id) || fixture.League.Round <= 5 {
		if fixture.League.Round <= 5 {
			tsid.Season--
		}
		tls := s.tlsService.GetTLS(tsid.GetTlsId())
		tsid = tls.GetTeamStatsId()
	}

	var hts, ats *model.TeamStats

	if mv, ok := s.statsMap[*tsid]; ok {
		hts = &mv // use the map value, since we have it
	} else {
		hts, _ = s.tsRepo.GetById(model.TeamStats{Id: *tsid})
	}


	return hts, ats
}

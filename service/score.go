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
	return s.getStatsTeam(fixture, true), s.getStatsTeam(fixture, false)
}

func (s *score) getStatsTeam(fixture *model.Fixture, home bool) *model.TeamStats {
	tsid := fixture.GetTeamStatsId(home)
	if s.statusService.IsScheduled(fixture.Fixture.Status.Id) || fixture.League.Round < 5 {
		if fixture.League.Round < 5 {
			tsid.Season--
		}
		tls := s.tlsService.GetTLS(tsid.GetTlsId())
		tsid = tls.GetTeamStatsId()
	}

	var ts *model.TeamStats

	// need to get by next fixture ID instead...
	if mv, ok := s.statsMap[*tsid]; ok {
		ts = &mv // use the map value, since we have it
	} else {
		ts, _ = s.tsRepo.GetById(model.TeamStats{Id: *tsid})
	}

	return ts
}
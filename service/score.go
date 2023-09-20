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
	prevStatsMap map[model.TeamStatsId]model.TeamStats
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
		prevStatsMap: make(map[model.TeamStatsId]model.TeamStats),
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
	id := *fixture.GetTeamStatsId(home)
	id.FixtureId = 0
	tsid := model.TeamStats{Id: id, NextFixtureId: fixture.Fixture.Id}

	if s.statusService.IsScheduled(fixture.Fixture.Status.Id) || fixture.League.Round < 5 {
		if fixture.League.Round < 5 {
			tsid.Id.Season--
		}
		tls := s.tlsService.GetTLS(tsid.Id.GetTlsId())
		tsid = model.TeamStats{Id: *tls.GetTeamStatsId(), NextFixtureId: 0}
	}

	// should use the maps first (prevStats by default, statsMap when using TLS)
	ts, _ := s.tsRepo.GetById(tsid)
	return ts
}
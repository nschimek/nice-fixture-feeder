package service

import (
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
)

type score struct {
	tsRepo repository.TeamStats
	tlsRepo repository.TeamLeagueSeason
	statusService FixtureStatus
	tlsMap map[model.TeamLeagueSeasonId]model.TeamLeagueSeason
	statsMap map[model.TeamStatsId]model.TeamStats
	scoreMap map[scores.ScoreId]scores.Score
}

func NewScore(tsRepo repository.TeamStats, 
	tlsRepo repository.TeamLeagueSeason,
	statusService FixtureStatus) *score {
	return &score{
		tsRepo: tsRepo, 
		tlsRepo: tlsRepo, 
		statusService: statusService,
		tlsMap: make(map[model.TeamLeagueSeasonId]model.TeamLeagueSeason),
		statsMap: make(map[model.TeamStatsId]model.TeamStats),
		scoreMap: make(map[scores.ScoreId]scores.Score),
	}
}

func (s *score) buildScoreMap() {
	s.scoreMap[scores.PointsStrengthId] = scores.NewPointsStrength(s.getStats)
}

func (s *score) getStats(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats) {
	var htsid, atsid *model.TeamStatsId

	// TODO: 1) if round is LTE 5, use previous season; 2) if fixture is NS, use TLS to get fixture ID
	htsid = fixture.GetTeamStatsId(true)
	atsid = fixture.GetTeamStatsId(false)

	var hts, ats *model.TeamStats

	if v, ok := s.statsMap[*htsid]; ok {
		hts = &v
	}

	if v, ok := s.statsMap[*atsid]; ok {
		ats = &v
	}

	return hts, ats
}

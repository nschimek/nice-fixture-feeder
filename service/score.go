package service

import (
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
)

type score struct {
	scores *scores.ScoreRegistry
	statusService FixtureStatus
	statsService *teamStats
	fixtureScores []model.FixtureScore
}

func NewScore(scores *scores.ScoreRegistry,
	statsService *teamStats, 
	statusService FixtureStatus) *score {
	s := &score{
		scores: scores,
		statusService: statusService,
		statsService: statsService,
	}
	s.setup()
	return s
}

func (s *score) setup() {
	s.scores.PointsStrength.SetStatsFunc(s.getStats)
}

func (s *score) ScoreAll(fixture *model.Fixture) {
	// score registry keeps a list of AllScores for use here
	for _, score := range s.scores.AllScores {
		if (score.CanScore(fixture)) {
			fs := score.Score(fixture)
			if fs != nil {
				s.fixtureScores = append(s.fixtureScores, *fs)
			}
		}
	}

	// TODO: next step is to persist!
}

func (s *score) getStats(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats, error) {	
	homeStats, err := s.getStatsTeam(fixture, true)
	if err != nil {
		return nil, nil, err
	}
	awayStats, err := s.getStatsTeam(fixture, false)
	if err != nil {
		return nil, nil, err
	}
	return homeStats, awayStats, nil
}

func (s *score) getStatsTeam(fixture *model.Fixture, home bool) (*model.TeamStats, error) {
	tsid := *fixture.GetTeamStatsNextId(home)
	
	// if the match has not been played or we are within the first 3 fixtures of the new season, we want to use TLS
	if s.statusService.IsScheduled(fixture.Fixture.Status.Id) || fixture.League.Round <= 3 {
		curr := false // by default, use TLS and then look-up stats by next fixture ID
		if fixture.League.Round <= 3 {
			tsid.Season--
			// the exception is when using last season's stats - in this case, we want the current
			curr = true 
		}
		return s.statsService.GetByIdWithTLS(tsid, curr)
	} 

	return s.statsService.GetById(tsid)
}
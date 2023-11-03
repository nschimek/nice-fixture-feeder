package service

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/nschimek/nice-fixture-feeder/repository"
	"github.com/nschimek/nice-fixture-feeder/service/scores"
)

//go:generate mockery --name Scoring --filename scoring_mock.go
type Scoring interface {
	SetFixtures(fixtures map[int]model.Fixture)
	AddFixturesFromMinMap(map[model.TeamLeagueSeasonId]int)
	Score()
}

type scoring struct {
	scores *scores.ScoreRegistry
	fixtureRepo repository.Fixture
	fixtureScoreRepo *repository.FixtureScore
	statusService FixtureStatus
	statsService TeamStats
	fixturesMap map[int]model.Fixture
	fixtureScores []model.FixtureScore
}

func NewScoring(scores *scores.ScoreRegistry,
	fixtureRepo repository.Fixture,
	fixtureScoreRepo *repository.FixtureScore,
	statsService TeamStats, 
	statusService FixtureStatus) *scoring {
	s := &scoring{
		scores: scores,
		fixtureRepo: fixtureRepo,
		fixtureScoreRepo: fixtureScoreRepo,
		statusService: statusService,
		statsService: statsService,
	}
	s.setup()
	return s
}

// Enable global access to set the Fixture Map
func (s *scoring) SetFixtures(fixtures map[int]model.Fixture) {
	s.fixturesMap = fixtures
}

// Add fixtures to be scored based on a Map of TLS ID by minimum ID
// The min ID, along with all future fixtures with the same league, season, and team will be scored.
// This will also exclude any Fixture IDs already in the Fixture map.
func (s *scoring) AddFixturesFromMinMap(fixtureMap map[model.TeamLeagueSeasonId]int) {
	notIn := core.MapToKeyArray[int, model.Fixture](s.fixturesMap)
	for tlsId, minId := range fixtureMap {
		fixtures, err := s.fixtureRepo.GetFutureFixturesByTLS(tlsId, minId, notIn)
		if err == nil {
			s.addFixtures(fixtures)
		}
	}
}

// Add fixtures for scoring.  Fixtures added here will not be re-queried by the above method.
func (s *scoring) addFixtures(fixtures []model.Fixture) {
	for _, f := range fixtures {
		s.fixturesMap[f.Fixture.Id] = f
	}
}


func (s *scoring) Score() {
	core.Log.WithField("fixtures", len(s.fixturesMap)).Info("Scoring Fixtures...")
	for _, fixture := range s.fixturesMap {
		s.scoreFixture(&fixture)
	}
	// Persist!
	core.Log.WithField("fixture_scores", len(s.fixtureScores)).Info("Persisting Fixture Scores...")
	s.fixtureScoreRepo.Upsert(s.fixtureScores)
}

func (s *scoring) setup() {
	s.scores.PointsStrength.SetStatsFunc(s.getStats)
}

func (s *scoring) scoreFixture(fixture *model.Fixture) {
	// score registry keeps a list of AllScores for use here
	for _, score := range s.scores.AllScores {
		if (score.CanScore(fixture)) {
			fs, err := score.Score(fixture)
			if fs != nil && err == nil {
				s.fixtureScores = append(s.fixtureScores, *fs)
			} else if err != nil {
				core.Log.WithField("fixture_id", fixture.Fixture.Id).Warn("Could not score fixture: ", err)
			}
		}
	}
}

func (s *scoring) getStats(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats, error) {	
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

func (s *scoring) getStatsTeam(fixture *model.Fixture, home bool) (*model.TeamStats, error) {
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
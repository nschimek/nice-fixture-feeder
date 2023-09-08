package scores

import (
	"math"

	"github.com/nschimek/nice-fixture-feeder/model"
)

type PointsStrength interface {
	Score
	SetStatsFunc(func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats))
}

type pointsStrength struct {
	StatsFunc func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats)
}

func NewPointsStrength() *pointsStrength {
	return &pointsStrength{}
}

func (s *pointsStrength) GetId() ScoreId {
	return PointsStrengthId
}

// we can only score if round was correcdtly parsed corectly
func (s *pointsStrength) CanScore(fixture *model.Fixture) bool {
	return fixture.League.Round > 0
}

func (s *pointsStrength) Score(fixture *model.Fixture) *model.FixtureScore {
	hs, as := s.StatsFunc(fixture)

	if hs == nil || as == nil {
		return nil
	}

	// Each team contributes their percentage of points equally.  If both teams have 100% of available points, it will be a perfect score.
	value := math.Round((s.pointsPercent(hs.Points, hs.TeamStatsFixtures.FixturesPlayed.Total) / 2) + 
		(s.pointsPercent(as.Points, as.TeamStatsFixtures.FixturesPlayed.Total) / 2)) * 100

	return &model.FixtureScore{
		FixtureId: fixture.Fixture.Id,
		ScoreId:   int(s.GetId()),
		Value: int(value),
	}
}

// Divides a team's current points total by the total number of points possible given the number of games they have played
func (s *pointsStrength) pointsPercent(points, played int) float64 {
	return math.Round((float64)(points / (played * 3)))
}

func (s *pointsStrength) SetStatsFunc(f func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats)) {
	s.StatsFunc = f
}
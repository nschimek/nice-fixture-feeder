package scores

import (
	"math"

	"github.com/nschimek/nice-fixture-feeder/model"
)

type PointsStrength struct {
	statsFunc func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats)
}

func NewPointsStrength(statsFunc func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats)) *PointsStrength {
	return &PointsStrength{
		statsFunc: statsFunc,
	}
}

func (s *PointsStrength) GetId() ScoreId {
	return PointsStrengthId
}

// we can only score if round was correcdtly parsed and we can get stats
func (s *PointsStrength) CanScore(fixture *model.Fixture) bool {
	hs, as := s.statsFunc(fixture)
	return fixture.League.Round > 0 && hs != nil && as != nil
}

func (s *PointsStrength) Score(fixture *model.Fixture) *model.FixtureScore {
	hs, as := s.statsFunc(fixture)

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
func (s *PointsStrength) pointsPercent(points, played int) float64 {
	return math.Round((float64)(points / (played * 3)))
}
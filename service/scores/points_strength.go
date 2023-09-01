package scores

import (
	"math"

	"github.com/nschimek/nice-fixture-feeder/model"
)

type PointsStrength struct {
	statsMap map[model.TeamStatsId]model.TeamStats
}

func NewPointsStrength(statsMap map[model.TeamStatsId]model.TeamStats) *PointsStrength {
	return &PointsStrength{
		statsMap: statsMap,
	}
}

func (s *PointsStrength) GetId() int {
	return int(PointsStrengthId)
}

func (s *PointsStrength) CanScore(fixture *model.Fixture) bool {
	return true
}

func (s *PointsStrength) Score(fixture *model.Fixture) *model.FixtureScore {
	hs := s.statsMap[*fixture.GetTeamStatsId(true)]
	as := s.statsMap[*fixture.GetTeamStatsId(false)]

	// Each team contributes their percentage of points equally.  If both teams have 100% of available points, it will be a perfect score.
	value := math.Round((s.pointsPercent(hs.Points, hs.TeamStatsFixtures.FixturesPlayed.Total) / 2) + 
		(s.pointsPercent(as.Points, as.TeamStatsFixtures.FixturesPlayed.Total) / 2)) * 100

	return &model.FixtureScore{
		FixtureId: fixture.Fixture.Id,
		ScoreId:   s.GetId(),
		Value: int(value),
	}
}

// Divides a team's current points total by the total number of points possible given the number of games they have played
func (s *PointsStrength) pointsPercent(points, played int) float64 {
	return math.Round((float64)(points / (played * 3)))
}
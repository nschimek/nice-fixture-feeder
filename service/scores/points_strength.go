package scores

import (
	"math"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/sirupsen/logrus"
)

type PointsStrength interface {
	Score
	SetStatsFunc(func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats, error))
}

type pointsStrength struct {
	StatsFunc func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats, error)
}

func NewPointsStrength() *pointsStrength {
	return &pointsStrength{}
}

func (s *pointsStrength) GetId() ScoreId {
	return PointsStrengthId
}

// we can only score if round was parsed corectly
func (s *pointsStrength) CanScore(fixture *model.Fixture) bool {
	return fixture.League.Round > 0
}

func (s *pointsStrength) Score(fixture *model.Fixture) (*model.FixtureScore, error) {
	hs, as, err := s.StatsFunc(fixture)

	if err != nil {
		return nil, err
	}

	// Each team contributes their percentage of points equally.  If both teams have 100% of available points, it will be a perfect score.
	value := math.Round((s.pointsPercent(hs.Points, hs.TeamStatsFixtures.FixturesPlayed.Total) / 2) + 
		(s.pointsPercent(as.Points, as.TeamStatsFixtures.FixturesPlayed.Total) / 2))

	core.Log.WithFields(logrus.Fields{
		"fixture_id": fixture.Fixture.Id,
		"hs_points": hs.Points,
		"hs_played": hs.TeamStatsFixtures.FixturesPlayed.Total,
		"as_points": as.Points,
		"as_played": as.TeamStatsFixtures.FixturesPlayed.Total,
		"value": value,
	}).Debug("Calculating Points Strength Score...")

	return &model.FixtureScore{
		FixtureId: fixture.Fixture.Id,
		ScoreId:   int(s.GetId()),
		Value: int(value),
	}, nil
}

// Divides a team's current points total by the total number of points possible given the number of games they have played
func (s *pointsStrength) pointsPercent(points, played int) float64 {
	return math.Round((float64(points) / (float64(played) * 3)) * 100)
}

func (s *pointsStrength) SetStatsFunc(f func(fixture *model.Fixture) (*model.TeamStats, *model.TeamStats, error)) {
	s.StatsFunc = f
}
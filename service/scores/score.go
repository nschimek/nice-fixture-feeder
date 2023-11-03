package scores

import "github.com/nschimek/nice-fixture-feeder/model"

type ScoreId int

const (
	PointsStrengthId ScoreId = 1
	PointsParityId ScoreId = 2
	GoalsId ScoreId = 3
	FormId ScoreId = 4
)

type Score interface {
	GetId() ScoreId
	CanScore(fixture *model.Fixture) bool
	Score(fixture *model.Fixture) (*model.FixtureScore, error)
}

type ScoreRegistry struct {
	AllScores []Score
	PointsStrength PointsStrength
}

func Setup() *ScoreRegistry {
	registry := &ScoreRegistry{
		AllScores: []Score{},
		PointsStrength: NewPointsStrength(),
	}

	registry.AllScores = append(registry.AllScores, registry.PointsStrength)

	return registry
}
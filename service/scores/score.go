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
	Score(fixture *model.Fixture) *model.FixtureScore
}

type ScoreRegistry struct {
	PointsStrength PointsStrength
}

func Setup() *ScoreRegistry {
	return &ScoreRegistry{
		PointsStrength: NewPointsStrength(),
	}
}
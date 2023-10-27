package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type FixtureScore struct {
	UpsertRepository[model.FixtureScore]
}

func NewFixtureScore(db core.Database) *FixtureScore {
	return &FixtureScore{
		UpsertRepository: upsertRepository[model.FixtureScore]{repository: newRepo(db, "fixture_scores")},
	}
}

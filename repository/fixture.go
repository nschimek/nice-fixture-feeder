package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type FixtureRepository struct {
	UpsertRepository[model.Fixture]
}

func NewFixtureRepository(db core.Database) *FixtureRepository {
	return &FixtureRepository{
		UpsertRepository: upsertRepository[model.Fixture]{repository: newRepo(db, "fixtures")},
	}
}
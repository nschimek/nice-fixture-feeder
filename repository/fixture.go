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
		UpsertRepository: upsertRepository[model.Fixture]{
			db: db,
			label: "fixtures",
			statsFunc: func(e []model.Fixture, r core.DatabaseResult, rs *ResultStats) {
				if r.Error == nil {
					rs.Success["fixture"] = len(e)
				} else {
					rs.Error["fixture"] = len(e)
				}
			},
		},
	}
}
package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type Fixture struct {
	UpsertRepository[model.Fixture]
}

func NewFixture(db core.Database) *Fixture {
	return &Fixture{
		UpsertRepository: upsertRepository[model.Fixture]{repository: newRepo(db, "fixtures")},
	}
}
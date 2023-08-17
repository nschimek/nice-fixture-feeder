package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type League struct {
	upsertRepository[model.League]
}

func NewLeague(db core.Database) *League {
	return &League{
		upsertRepository: upsertRepository[model.League]{
			repository: newRepo(db, "leagues and seasons"),
		},
	}
}
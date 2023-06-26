package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type LeagueRepository struct {
	upsertRepository[model.League]
}

func NewLeagueRepository(db core.Database) *LeagueRepository {
	return &LeagueRepository{
		upsertRepository: upsertRepository[model.League]{
			repository: newRepo(db, "leagues and seasons"),
		},
	}
}
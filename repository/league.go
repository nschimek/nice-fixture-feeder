package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type LeagueRepository struct {
	repository[model.League]
}

func NewLeagueRepository(db core.Database) *LeagueRepository {
	return &LeagueRepository{
		repository: repository[model.League]{
			DB: db,
			label: "leagues",
			statsFunc: func(e []model.League, r core.DatabaseResult, rs *ResultStats) {
				if r.Error == nil {
					rs.Success["league"] = len(e)
					rs.Success["season"] = len(e)
				} else {
					rs.Error["league"] = len(e)
					rs.Error["season"] = len(e)
				}
			},
		},
	}
}
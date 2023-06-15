package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type TeamStatsRepository interface {
	UpsertRepository[model.TeamStats]
	GetByIdRepository[model.TeamStats, model.TeamStats]
}

type teamStatsRepository struct {
	upsertRepository[model.TeamStats]
	getByIdRepository[model.TeamStats, model.TeamStats]
}

func NewTeamStatsRepository(db core.Database) *teamStatsRepository {
	return &teamStatsRepository{
		upsertRepository: upsertRepository[model.TeamStats]{
			DB: db,
			label: "team_stats",
			statsFunc: func(e []model.TeamStats, r core.DatabaseResult, rs *ResultStats) {
				if r.Error == nil {
					rs.Success["team_stats"] = len(e)
				} else {
					rs.Error["team_stats"] = len(e)
				}
			},
		},
		getByIdRepository: getByIdRepository[model.TeamStats, model.TeamStats]{db: db},
	}
}
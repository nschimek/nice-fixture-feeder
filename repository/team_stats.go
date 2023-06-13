package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type TeamStatsRepository struct {
	upsertRepository[model.TeamStats]
}

func NewTeamStatsRepository(db core.Database) *TeamStatsRepository {
	return &TeamStatsRepository{
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
	}
}
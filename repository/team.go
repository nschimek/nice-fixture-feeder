package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type TeamRepository struct {
	repository[model.Team]
}

func NewTeamRepository(db core.Database) *TeamRepository {
	return &TeamRepository{
		repository: repository[model.Team]{
			DB: db,
			label: "teams",
			statsFunc: func(e []model.Team, r core.DatabaseResult, rs *ResultStats) {
				if r.Error == nil {
					rs.Success["team"] = len(e)
					rs.Success["team_league_season"] = len(e)
				} else {
					rs.Error["team"] = len(e)
					rs.Error["team_league_season"] = len(e)
				}
			},
		},
	}
}
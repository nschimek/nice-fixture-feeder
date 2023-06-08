package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type LeagueRepository struct {
	DB core.Database
}

func (lr *LeagueRepository) Upsert(leagues []model.League) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField("leagues", len(leagues)).Infof("Create/Updating Leagues and League Seasons...")

	r := lr.DB.Upsert(&leagues)

	if r.Error == nil {
		rs.Success["team"] = len(leagues)
		rs.Success["team_league_season"] = len(leagues)
		core.Log.WithField("teams", len(leagues)).Infof("Successfully create/updated teams along with team league seasons")
	} else {
		rs.Error["team"] = len(leagues)
		rs.Error["team_league_season"] = len(leagues)
	}

	return rs
}
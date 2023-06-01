package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type TeamRepository struct {
	DB core.Database
}

func (lr *TeamRepository) Upsert(teams []model.Team) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField("teams", len(teams)).Infof("Create/Updating Leagues and Team League Seasons...")

	r := lr.DB.Upsert(&teams)

	if r.Error == nil {
		rs.Success["team"] = len(teams)
		rs.Success["team_league_season"] = len(teams)
		core.Log.WithField("teams", len(teams)).Infof("Successfully create/updated teams along with team league seasons")
	} else {
		rs.Error["team"] = len(teams)
		rs.Error["team_league_season"] = len(teams)
	}

	return rs
}
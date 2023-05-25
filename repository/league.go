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

	for _, league := range leagues {
		r := lr.DB.Upsert(&league)

		if r.Error == nil {
			rs.Success["league"]++
			rs.Success["season"] = rs.Success["season"] + len(league.Seasons)
			core.Log.WithField("seasons", len(league.Seasons)).Infof("Successfully create/updated league %d along with seasons", league.League.Id)
		} else {
			rs.Error["league"]++
			rs.Error["season"] = rs.Error["season"] + len(league.Seasons)
		}
	}

	return rs
}
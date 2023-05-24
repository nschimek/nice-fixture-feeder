package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type LeagueRepository struct {
	DB core.Database
}

func (lr *LeagueRepository) UpsertLeagues(leagues []model.League) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField("leagues", len(leagues)).Infof("Create/Updating Leagues and League Seasons...")

	for i, league := range leagues {
		r1 := lr.DB.UpsertWithOmit(&league, "Seasons")

		if r1.Error == nil && len(league.Seasons) > 0 {
			rs.Success["league"]++
			for i := range league.Seasons {
				league.Seasons[i].LeagueId = league.League.Id
			}
			r2 := lr.DB.Upsert(league.Seasons)
			if r2.Error == nil {
				core.Log.WithField("seasons", len(league.Seasons)).Infof("Successfully create/updated league %d along with seasons", league.League.Id)
				rs.Success["season"] = rs.Success["season"] + len(league.Seasons)
			} else {
				rs.Error["season"] = rs.Error["season"] + len(league.Seasons)
			}
		} else if r1.Error != nil {
			rs.Error["league"]++
			l := &leagues[i] // have to get the pointer to change it
			l.CaptureError(r1.Error)
		}
	}

	return rs
}
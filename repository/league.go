package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"gorm.io/gorm/clause"
)

var updateAll = clause.OnConflict{UpdateAll: true}

type LeagueRepository struct {
	DB *core.Database
}

func (lr *LeagueRepository) UpsertLeagues(leagues []model.League) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField("leagues", len(leagues)).Infof("Create/Updating Leagues and League Seasons...")

	for _, league := range leagues {
		r1 := lr.DB.Gorm.Clauses(updateAll).Omit("Seasons").Create(&league)

		if r1.Error == nil && len(league.Seasons) > 0 {
			rs.Success["season"]++
			for i := range league.Seasons {
				league.Seasons[i].LeagueId = league.League.Id
			}
			r2 := lr.DB.Gorm.Clauses(updateAll).Create(league.Seasons)
			if r2.Error == nil {
				core.Log.WithField("seasons", len(league.Seasons)).Infof("Successfully create/updated league %d along with seasons", league.League.Id)
				rs.Success["league"]++
			} else {
				rs.Error["league"]++
			}
		} else if r1.Error != nil {
			rs.Error["season"]++
			league.CaptureError(r1.Error)
		}
	}

	return rs
}
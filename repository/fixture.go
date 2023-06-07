package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type FixtureRepository struct {
	DB core.Database
}

func (fr *FixtureRepository) Upsert(fixtures []model.Fixture) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField("fixtures", len(fixtures)).Infof("Create/Updating Fixtures...")

	r := fr.DB.Upsert(fixtures)

	if r.Error == nil {
		rs.Success["fixture"] = len(fixtures)
		core.Log.WithField("fixtures", len(fixtures)).Infof("Successfully create/updated fixtures")
	} else {
		rs.Error["fixture"] = len(fixtures)
	}

	return rs
}
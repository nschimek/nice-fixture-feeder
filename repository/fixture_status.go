package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

type FixtureStatusRepository struct {
	DB core.Database
}

func (r *FixtureStatusRepository) GetById(id string) *model.FixtureStatus {
	var status model.FixtureStatus
	r.DB.GetById(id, &status)
	return &status
}
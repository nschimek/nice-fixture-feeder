package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name FixtureStatusRepository --filename fixture_status_mock.go
type FixtureStatusRepository interface {
	GetAll() []model.FixtureStatus
}

type fixtureStatusRepository struct {
	db core.Database
}

func NewFixtureStatusRepository(db core.Database) FixtureStatusRepository {
	return &fixtureStatusRepository{db: db}
}

func (r *fixtureStatusRepository) GetAll() []model.FixtureStatus {
	var statues []model.FixtureStatus
	r.db.GetAll(&statues)
	return statues
}
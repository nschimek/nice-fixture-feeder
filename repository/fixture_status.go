package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name FixtureStatusRepository --filename fixture_status_mock.go
type FixtureStatusRepository interface {
	GetAll() ([]model.FixtureStatus, error)
}

type fixtureStatusRepository struct {
	db core.Database
}

func NewFixtureStatusRepository(db core.Database) FixtureStatusRepository {
	return &fixtureStatusRepository{db: db}
}

func (r *fixtureStatusRepository) GetAll() ([]model.FixtureStatus, error) {
	var statues []model.FixtureStatus
	if err := r.db.GetAll(&statues).Error; err != nil {
		return nil, err
	}
	return statues, nil
}
package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
)

//go:generate mockery --name FixtureStatus --filename fixture_status_mock.go
type FixtureStatus interface {
	GetAll() ([]model.FixtureStatus, error)
}

type fixtureStatus struct {
	db core.Database
}

func NewFixtureStatus(db core.Database) FixtureStatus {
	return &fixtureStatus{db: db}
}

func (r *fixtureStatus) GetAll() ([]model.FixtureStatus, error) {
	var statues []model.FixtureStatus
	if err := r.db.GetAll(&statues).Error; err != nil {
		return nil, err
	}
	return statues, nil
}
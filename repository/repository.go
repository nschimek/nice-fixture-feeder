package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
)

var Repositories *RepositoryRegistry

type RepositoryRegistry struct {
	Fixture *FixtureRepository
	FixtureStatus FixtureStatusRepository
	League *LeagueRepository
	TeamLeagueSeason TeamLeagueSeasonRepository
	TeamStats TeamStatsRepository
	Team *TeamRepository
}

func Setup(db core.Database) {
	Repositories = &RepositoryRegistry{
		Fixture: NewFixtureRepository(db),
		FixtureStatus: NewFixtureStatusRepository(db),
		League: NewLeagueRepository(db),
		TeamLeagueSeason: NewTeamLeagueSeasonRepository(db),
		TeamStats: NewTeamStatsRepository(db),
		Team: NewTeamRepository(db),
	}
}

type repository struct {
	db core.Database
	label string
}

func newRepo(db core.Database, label string) *repository {
	return &repository{db: db, label: label}
}

//go:generate mockery --name UpsertRepository --filename repository_upsert_mock.go
type UpsertRepository[T any] interface {
	Upsert(entities []T) ([]T, error)
}

type upsertRepository[T any] struct {
	*repository
}

func (r upsertRepository[T]) Upsert(entities []T) ([]T, error) {
	core.Log.WithField(r.label, len(entities)).Infof("Create/updating %s...", r.label)

	if len(entities) == 0 || entities == nil {
		core.Log.Warnf("Got no or nil %s, skipping persistence!", r.label)
		return nil, nil
	}

	res := r.db.Upsert(&entities)

	if res.Error != nil {
		core.Log.WithField(r.label, len(entities)).Error("Issues during persistence")
		return nil, res.Error
	}

	core.Log.WithField(r.label, len(entities)).Info("Persistence successful!")
	return entities, nil
}

//go:generate mockery --name GetByIdRepository --filename repository_id_mock.go
type GetByIdRepository[T any, I any] interface {
	GetById(id I) (*T, error)
}

type getByIdRepository[T any, I any] struct {
	*repository
}

func (r getByIdRepository[T, I]) GetById(id I) (*T, error) {
	var dest T
	if res := r.db.GetById(id, &dest); res.RowsAffected == 0 && res.Error == nil {
		return nil, nil
	} else if res.Error != nil {
		return nil, res.Error
	}
	return &dest, nil
}

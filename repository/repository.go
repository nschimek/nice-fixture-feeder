package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
)

type RepositoryRegistry struct {
	Fixture Fixture
	FixtureStatus FixtureStatus
	FixtureScore *FixtureScore
	League *League
	TeamLeagueSeason TeamLeagueSeason
	TeamStats TeamStats
	Team *Team
}

// Setup the Repositories so that they can be injected in Services and Requesters
func Setup(db core.Database) *RepositoryRegistry {
	return &RepositoryRegistry{
		Fixture: NewFixture(db),
		FixtureStatus: NewFixtureStatus(db),
		FixtureScore: NewFixtureScore(db),
		League: NewLeague(db),
		TeamLeagueSeason: NewTeamLeagueSeason(db),
		TeamStats: NewTeamStats(db),
		Team: NewTeam(db),
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
	UpsertOne(entity T) (T, error)
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

// TODO: We should probably return a pointer to T instead of just T here.
func (r upsertRepository[T]) UpsertOne(entity T) (T, error) {
	res, err := r.Upsert([]T{entity})
	if err != nil {
		return *new(T), err
	}
	return res[0], err 
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

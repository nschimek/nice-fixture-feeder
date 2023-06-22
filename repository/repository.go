package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
)

//go:generate mockery --name UpsertRepository --filename repository_upsert_mock.go
type UpsertRepository[T any] interface {
	Upsert(entities []T) *ResultStats
}

type upsertRepository[T any] struct {
	db core.Database
	label string
	statsFunc func(e []T, r core.DatabaseResult, rs *ResultStats)
}

func (u upsertRepository[T]) Upsert(entities []T) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField(u.label, len(entities)).Infof("Create/updating %s...", u.label)

	res := u.db.Upsert(&entities)

	if (u.statsFunc != nil) {
		u.statsFunc(entities, res, rs)
	} else {
		u.defaultStatsFunc(entities, res, rs)
	}

	return rs
}

func (u upsertRepository[T]) defaultStatsFunc(e []T, r core.DatabaseResult, rs *ResultStats) {
	if r.Error == nil {
		rs.Success[u.label] = len(e)
	} else {
		rs.Error[u.label] = len(e)
	}
}

//go:generate mockery --name GetByIdRepository --filename repository_id_mock.go
type GetByIdRepository[T any, I any] interface {
	GetById(id I) *T
}

type getByIdRepository[T any, I any] struct {
	db core.Database
}

func (r getByIdRepository[T, I]) GetById(id I) *T {
	var dest T
	if c := r.db.GetById(id, &dest).RowsAffected; c == 0 {
		return nil
	}
	return &dest
}

//go:generate mockery --name SaveRepository --filename repository_save_mock.go
type SaveRepository[T any] interface {
	Save(entity *T) (*T, error)
}

type saveRepository[T any] struct {
	db core.Database
}

func (r saveRepository[T]) Save(entity *T) (*T, error) {
	if err := r.db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

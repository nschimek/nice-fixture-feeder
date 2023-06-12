package repository

import "github.com/nschimek/nice-fixture-feeder/core"

//go:generate mockery --name Repository
type Repository[T any] interface {
	Upsert([]T) *ResultStats
}

type repository[T any] struct {
	DB core.Database
	label string
	statsFunc func(e []T, r core.DatabaseResult, rs *ResultStats)
}

func (r *repository[T]) Upsert(entities []T) *ResultStats {
	rs := NewResultStats()
	core.Log.WithField(r.label, len(entities)).Infof("Create/updating %s...", r.label)

	res := r.DB.Upsert(&entities)

	r.statsFunc(entities, res, rs)

	return rs
}
package repository

//go:generate mockery --name Repository
type Repository[T any] interface {
	Upsert([]T) *ResultStats
}
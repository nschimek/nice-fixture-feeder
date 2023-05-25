package repository

type Repository[T any] interface {
	Upsert([]T) *ResultStats
}
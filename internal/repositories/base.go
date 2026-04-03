package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type BaseRepository[T any] struct {
	db *pgxpool.Pool
}

func NewBaseRepository[T any](db *pgxpool.Pool) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

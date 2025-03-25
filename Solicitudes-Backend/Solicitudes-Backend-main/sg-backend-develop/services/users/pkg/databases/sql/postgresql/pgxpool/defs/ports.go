package defs

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Connect(Config) error
	Close()
	Pool() *pgxpool.Pool
}

type Config interface {
	Validate() error
	DNS() string
}

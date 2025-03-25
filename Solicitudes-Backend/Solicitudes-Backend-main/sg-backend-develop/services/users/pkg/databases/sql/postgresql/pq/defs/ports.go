package defs

import (
	"database/sql"
)

type Repository interface {
	Connect(Config) error
	Close()
	DB() *sql.DB
}

type Config interface {
	Validate() error
	DNS() string
}

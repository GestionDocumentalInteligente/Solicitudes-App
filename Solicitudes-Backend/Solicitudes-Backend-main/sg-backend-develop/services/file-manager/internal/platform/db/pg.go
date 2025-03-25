package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/teamcubation/sg-file-manager-api/internal/platform/env"
)

var db *sqlx.DB //nolint:gochecknoglobals

func NewPGConnection(ctx context.Context) (*sqlx.DB, error) {
	if db == nil || db.PingContext(ctx) != nil {
		conn, err := dbConnectionURL()
		if err != nil {
			return nil, err
		}

		db, err = sqlx.ConnectContext(ctx, "postgres", conn)
		if err != nil {
			log.Println("########## DB ERROR: " + err.Error() + " #############")
			return nil, fmt.Errorf("### DB ERROR: %w", err)
		}
	}

	return db, nil
}

func CloseDB() error {
	if db != nil {
		err := db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func dbConnectionURL() (string, error) {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require",
		env.GetDBUser(),
		env.GetDBPass(),
		env.GetDBHost(),
		env.GetDBName()), nil
}

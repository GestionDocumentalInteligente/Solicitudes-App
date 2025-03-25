package sdkpostgresql

import (
	"os"

	"github.com/teamcubation/sg-backend/pkg/databases/sql/postgresql/pgxpool/defs"
)

// NOTE: Diseñado para establer conexion con 1 base de datos durante la ejecución de la app
func Bootstrap() (defs.Repository, error) {
	config := newConfig(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_MIGRATIONS_DIR"),
		os.Getenv("POSTGRES_DB"),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

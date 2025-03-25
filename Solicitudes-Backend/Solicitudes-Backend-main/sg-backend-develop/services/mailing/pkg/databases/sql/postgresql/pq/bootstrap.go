package sdkpg

import (
	"os"

	defs "github.com/teamcubation/sg-mailing-api/pkg/databases/sql/postgresql/pq/defs"
)

func Bootstrap(dbNameKey string) (defs.Repository, error) {
	config := newConfig(
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv(dbNameKey),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

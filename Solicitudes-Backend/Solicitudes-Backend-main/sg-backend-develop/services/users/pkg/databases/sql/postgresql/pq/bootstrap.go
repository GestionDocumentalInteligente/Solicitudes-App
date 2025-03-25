package sdkpg

import (
	"github.com/spf13/viper"

	defs "github.com/teamcubation/sg-users/pkg/databases/sql/postgresql/pq/defs"
)

func Bootstrap(dbNameKey string) (defs.Repository, error) {
	config := newConfig(
		viper.GetString("POSTGRES_USERNAME"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("POSTGRES_HOST"),
		viper.GetString("POSTGRES_PORT"),
		viper.GetString(dbNameKey),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newRepository(config)
}

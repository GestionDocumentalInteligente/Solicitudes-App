package sdkjwt

import (
	"github.com/spf13/viper"

	"github.com/teamcubation/sg-backend/pkg/jwt/v5/defs"
)

func Bootstrap(secretKey, accessExpirationKey, refreshExpirationKey string) (defs.Service, error) {
	config := newConfig(
		viper.GetString(secretKey),
		viper.GetInt(accessExpirationKey),
		viper.GetInt(refreshExpirationKey),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newService(config)
}

package sdkjwt

import (
	"github.com/spf13/viper"

	"github.com/teamcubation/sg-users/pkg/jwt/v5/defs"
)

func Bootstrap(secretKey, expirationKey string) (defs.Service, error) {
	config := newConfig(
		viper.GetString(secretKey),
		viper.GetInt(expirationKey),
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newService(config)
}

package sdkjwt

import (
	"os"

	"github.com/teamcubation/sg-auth/pkg/jwt/v5/ports"
)

func Bootstrap() (ports.Service, error) {
	config, err := newConfig(
		os.Getenv("JWT_SECRET_KEY"),
	)
	if err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newService(config)
}

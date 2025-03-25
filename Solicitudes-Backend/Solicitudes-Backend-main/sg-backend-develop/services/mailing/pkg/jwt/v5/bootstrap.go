package sdkjwt

import (
	"os"

	"github.com/teamcubation/sg-mailing-api/pkg/jwt/v5/defs"
)

func Bootstrap(secretKey, expirationKey string) (defs.Service, error) {
	config := newConfig(
		os.Getenv(secretKey),
		60,
	)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return newService(config)
}

package sdkjwt

import (
	"fmt"
	"time"

	"github.com/teamcubation/sg-mailing-api/pkg/jwt/v5/defs"
)

type config struct {
	secret            string
	expirationMinutes int
}

func newConfig(secretKey string, expirationMinutes int) defs.Config {
	return &config{
		secret:            secretKey,
		expirationMinutes: expirationMinutes,
	}
}

func (c *config) GetSecretKey() string {
	return c.secret
}

func (c *config) GetExpirationMinutes() time.Duration {
	return time.Duration(c.expirationMinutes) * time.Minute
}

func (c *config) Validate() error {
	if c.secret == "" {
		return fmt.Errorf("JWT secret key is not configured")
	}
	if c.expirationMinutes <= 0 {
		return fmt.Errorf("JWT expiration minutes must be greater than zero")
	}
	return nil
}

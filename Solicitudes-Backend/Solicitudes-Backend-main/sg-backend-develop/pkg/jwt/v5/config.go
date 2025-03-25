package sdkjwt

import (
	"fmt"
	"time"

	"github.com/teamcubation/sg-backend/pkg/jwt/v5/defs"
)

type config struct {
	secret                   string
	accessExpirationMinutes  int
	refreshExpirationMinutes int
}

// newConfig crea una nueva configuración de JWT
func newConfig(secretKey string, accessExpirationMinutes, refreshExpirationMinutes int) defs.Config {
	return &config{
		secret:                   secretKey,
		accessExpirationMinutes:  accessExpirationMinutes,
		refreshExpirationMinutes: refreshExpirationMinutes,
	}
}

// GetSecretKey devuelve la clave secreta para firmar los tokens JWT
func (c *config) GetSecretKey() string {
	return c.secret
}

// GetAccessExpiration devuelve la duración de expiración del token de acceso
func (c *config) GetAccessExpiration() time.Duration {
	return time.Duration(c.accessExpirationMinutes) * time.Minute
}

// GetRefreshExpiration devuelve la duración de expiración del token de refresco
func (c *config) GetRefreshExpiration() time.Duration {
	return time.Duration(c.refreshExpirationMinutes) * time.Minute
}

// Validate verifica que la configuración de JWT sea válida
func (c *config) Validate() error {
	if c.secret == "" {
		return fmt.Errorf("JWT secret key is not configured")
	}
	if c.accessExpirationMinutes <= 0 {
		return fmt.Errorf("JWT access expiration minutes must be greater than zero")
	}
	if c.refreshExpirationMinutes <= 0 {
		return fmt.Errorf("JWT refresh expiration minutes must be greater than zero")
	}
	return nil
}

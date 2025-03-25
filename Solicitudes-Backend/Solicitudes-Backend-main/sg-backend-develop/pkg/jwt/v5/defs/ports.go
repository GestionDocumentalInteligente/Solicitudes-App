// defs/types.go

package defs

import (
	"context"
	"time"
)

// Config define la interfaz para la configuración del servicio JWT.
type Config interface {
	GetAccessExpiration() time.Duration
	GetRefreshExpiration() time.Duration
	GetSecretKey() string
	Validate() error
}

// Service define la interfaz para el servicio JWT.
type Service interface {
	GenerateTokens(context.Context, string) (string, string, error)
	ValidateToken(context.Context, string) (*TokenClaims, error)
	GetAccessExpiration() time.Duration
	GetRefreshExpiration() time.Duration
	ValidateTokenAllowExpired(ctx context.Context, tokenString string) (*TokenClaims, error)
}

// defs/types.go

package defs

import (
	"context"
	"time"
)

// Config define la interfaz para la configuración del servicio JWT.
type Config interface {
	GetExpirationMinutes() time.Duration
	GetSecretKey() string
	Validate() error
}

// Service define la interfaz para el servicio JWT.
type Service interface {
	GenerateToken(context.Context, string) (string, error)
	ValidateToken(context.Context, string) (*TokenClaims, error)
}

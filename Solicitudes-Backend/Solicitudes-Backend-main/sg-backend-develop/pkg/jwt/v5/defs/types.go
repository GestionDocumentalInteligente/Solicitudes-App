package defs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims representa las claims personalizadas para el token JWT.
type Claims struct {
	Subject string `json:"sub"`
	jwt.RegisteredClaims
}

// TokenClaims representa las claims extraídas de un token validado.
type TokenClaims struct {
	Subject   string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

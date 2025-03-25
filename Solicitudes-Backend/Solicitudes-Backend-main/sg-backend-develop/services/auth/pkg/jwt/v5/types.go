package sdkjwt

import (
	"time"
)

// Token representa un token de autenticación JWT
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type TokenClaims struct {
	Subject   string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

//type MapClaims jwt.MapClaims

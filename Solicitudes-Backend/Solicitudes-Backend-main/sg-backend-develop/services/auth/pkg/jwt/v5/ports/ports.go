package ports

import "github.com/golang-jwt/jwt/v5"

type Service interface {
	GenerateToken(claims jwt.MapClaims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type Config interface {
	GetSecretKey() string
	Validate() error
}

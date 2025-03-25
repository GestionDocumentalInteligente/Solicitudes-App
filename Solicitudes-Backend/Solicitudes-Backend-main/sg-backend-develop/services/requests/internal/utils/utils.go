package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetCuilFromJwt(c *gin.Context) (int64, error) {
	token, exists := c.Get("token")
	if !exists {
		return 0, errors.New("Token not found")
	}

	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return 0, errors.New("Invalid token type")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claims")
	}

	userID, ok := claims["user_id"].(float64) // JWT numéricos se manejan como float64
	if !ok {
		return 0, errors.New("User ID not found in token")
	}

	return int64(userID), nil
}

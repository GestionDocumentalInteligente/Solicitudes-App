package authhdl

import (
	"github.com/gin-gonic/gin"
)

type AuthHandlerRouter struct {
	authhdl *AuthHandler
}

func NewRouter(authhdl *AuthHandler) *AuthHandlerRouter {
	return &AuthHandlerRouter{
		authhdl: authhdl,
	}
}

func (s *AuthHandlerRouter) AddRoutesV1(v1 *gin.RouterGroup) {
	v1.POST("/login", s.authhdl.Login)
}

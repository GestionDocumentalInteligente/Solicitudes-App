package authhdl

import (
	"net/http"

	"github.com/teamcubation/sg-file-manager-api/internal/manager"
	"github.com/teamcubation/sg-file-manager-api/internal/manager/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCase manager.AuthUseCase
}

func NewAuthHandler(uc manager.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: uc,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials auth.Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.useCase.Login(c, credentials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse(result))
}

package transport

import (
	sdktypes "github.com/teamcubation/sg-auth/pkg/types"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func DomainToLoginResponse(lc *sdktypes.LoginCredentials) *sdktypes.LoginCredentials {
	return &sdktypes.LoginCredentials{
		Username:     lc.Username,
		PasswordHash: lc.PasswordHash,
	}
}

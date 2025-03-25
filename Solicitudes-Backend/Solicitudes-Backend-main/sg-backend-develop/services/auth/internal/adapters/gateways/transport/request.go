package transport

import sdktypes "github.com/teamcubation/sg-auth/pkg/types"

func LoginRequestToDomain(lr *sdktypes.LoginCredentials) *sdktypes.LoginCredentials {
	return &sdktypes.LoginCredentials{
		Username:     lr.Username,
		PasswordHash: lr.PasswordHash,
	}
}

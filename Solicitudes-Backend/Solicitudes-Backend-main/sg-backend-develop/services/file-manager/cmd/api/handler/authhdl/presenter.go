package authhdl

import (
	"github.com/teamcubation/sg-file-manager-api/internal/manager/auth"
)

type TokenResponseJSON struct {
	Token string `json:"token"`
}

func TokenResponse(i auth.TokenResponse) TokenResponseJSON {
	var response TokenResponseJSON

	response.Token = i.Token

	return response
}

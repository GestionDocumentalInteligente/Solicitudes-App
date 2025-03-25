package mailconn

import (
	sdlsmtpdefs "github.com/teamcubation/sg-mailing-api/pkg/mailing/smtp/defs"

	dto "github.com/teamcubation/sg-mailing-api/internal/core/dto"
)

func toSdkEmailData(data *dto.EmailData) *sdlsmtpdefs.EmailData {
	return &sdlsmtpdefs.EmailData{
		Email:        data.Email,
		Name:         data.Name,
		Subject:      data.Subject,
		BodyTemplate: data.BodyTemplate,
		Token:        data.Token,
	}
}

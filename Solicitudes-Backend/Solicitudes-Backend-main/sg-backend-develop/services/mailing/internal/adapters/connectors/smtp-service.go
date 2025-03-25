package mailconn

import (
	"context"
	"log"

	sdksmtp "github.com/teamcubation/sg-mailing-api/pkg/mailing/smtp"
	sdlsmtpdefs "github.com/teamcubation/sg-mailing-api/pkg/mailing/smtp/defs"

	dto "github.com/teamcubation/sg-mailing-api/internal/core/dto"
	ports "github.com/teamcubation/sg-mailing-api/internal/core/ports"
)

type SmtpService struct {
	smtpService sdlsmtpdefs.Service
}

func NewSmtpService() (ports.SmtpService, error) {
	smtpService, err := sdksmtp.Bootstrap()
	if err != nil {
		log.Fatalf("Failed to initialize SMTP service: %v", err)
	}

	return &SmtpService{
		smtpService: smtpService,
	}, nil
}

func (ss *SmtpService) SendVerificationEmail(ctx context.Context, data *dto.EmailData) error {
	return ss.smtpService.SendVerificationEmail(ctx, toSdkEmailData(data))
}

func (ss *SmtpService) SendNewRequestEmail(ctx context.Context, code string, data *dto.EmailData) error {
	return ss.smtpService.SendNewRequestEmail(ctx, code, toSdkEmailData(data))
}

func (ss *SmtpService) SendUpdateRequestByCodeMessage(ctx context.Context, code, obs string, data *dto.EmailData) error {
	return ss.smtpService.SendUpdateRequestByCodeMessage(ctx, code, obs, toSdkEmailData(data))
}

func (ss *SmtpService) SendUpdateRequestMessage(ctx context.Context, code string, data *dto.EmailData) error {
	return ss.smtpService.SendUpdateRequestMessage(ctx, code, toSdkEmailData(data))
}

func (ss *SmtpService) SendValidateRequestMessage(ctx context.Context, code string, data *dto.EmailData) error {
	return ss.smtpService.SendValidateRequestMessage(ctx, code, toSdkEmailData(data))
}

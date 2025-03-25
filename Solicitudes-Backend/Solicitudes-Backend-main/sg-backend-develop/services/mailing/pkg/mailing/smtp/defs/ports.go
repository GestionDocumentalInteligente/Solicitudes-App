package defs

import (
	"context"
	"net/smtp"
)

// Config define la interfaz que debe cumplir la configuración SMTP
type Config interface {
	GetSMTPServer() string
	GetAuth() smtp.Auth
	GetFrom() string
	GetPort() string
	GetVerificationURL() string
	Validate() error
}

type Service interface {
	SendVerificationEmail(context.Context, *EmailData) error
	SendNewRequestEmail(ctx context.Context, code string, data *EmailData) error
	SendUpdateRequestByCodeMessage(ctx context.Context, code, obs string, data *EmailData) error
	SendUpdateRequestMessage(ctx context.Context, code string, data *EmailData) error
	SendValidateRequestMessage(ctx context.Context, code string, data *EmailData) error
}

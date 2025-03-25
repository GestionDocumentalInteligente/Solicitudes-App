package mailing

import (
	"context"
	"fmt"

	dto "github.com/teamcubation/sg-mailing-api/internal/core/dto"
	"github.com/teamcubation/sg-mailing-api/internal/core/entities"
	ports "github.com/teamcubation/sg-mailing-api/internal/core/ports"
)

const subject = "Subject: Bienvenido/a solicitudes online San Isidro\r\n"
const mime = "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

type UseCases struct {
	jwtService  ports.JwtService
	smtpService ports.SmtpService
	repo        ports.Repository
}

func NewUseCases(js ports.JwtService, ss ports.SmtpService, repo ports.Repository) ports.UseCases {
	return &UseCases{
		smtpService: ss,
		jwtService:  js,
		repo:        repo,
	}
}

func (u *UseCases) InitiateEmailVerification(ctx context.Context, email, name string) error {
	data := &dto.EmailData{
		Email:        email,
		Name:         name,
		Subject:      subject,
		BodyTemplate: mime,
	}

	token, err := u.jwtService.GenerateToken(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	data.Token = token

	// Enviar el correo de verificación
	if err := u.smtpService.SendVerificationEmail(ctx, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (u *UseCases) ActivateAccount(ctx context.Context, token string) error {
	claims, err := u.jwtService.ValidateToken(ctx, token)
	if err != nil {
		return &entities.ErrorTokenInvalid{
			Msg: "token could not be validated or is invalid",
		}
	}

	return u.repo.UpdateUser(ctx, claims.Subject)
}

func (u *UseCases) ResendActivationEmail(ctx context.Context, token string) error {
	claims, err := u.jwtService.GetTokenInfo(ctx, token)
	if err != nil {
		return &entities.ErrorTokenInvalid{
			Msg: "token could not be validated or is invalid",
		}
	}

	user, err := u.repo.GetuserByEmail(ctx, claims.Subject)
	if err != nil {
		return fmt.Errorf("error in repository: %w", err)
	}

	if user == nil {
		return &entities.NotFoundInDatabase{
			Msg: "user not found",
		}
	}

	if user.EmailValidated {
		return &entities.UserAlreadyActive{
			Email: claims.Subject,
		}
	}

	data := &dto.EmailData{
		Email:        claims.Subject,
		Name:         claims.Subject,
		Subject:      subject,
		BodyTemplate: mime,
	}

	newToken, err := u.jwtService.GenerateToken(ctx, claims.Subject)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	data.Token = newToken

	if err := u.smtpService.SendVerificationEmail(ctx, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (u *UseCases) ResendActivationEmailExistingUser(ctx context.Context, email string) error {
	user, err := u.repo.GetuserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error in repository: %w", err)
	}

	if user == nil {
		return &entities.NotFoundInDatabase{
			Msg: "user not found",
		}
	}

	if user.EmailValidated {
		return &entities.UserAlreadyActive{
			Email: email,
		}
	}

	data := &dto.EmailData{
		Email:        email,
		Name:         email,
		Subject:      subject,
		BodyTemplate: mime,
	}

	newToken, err := u.jwtService.GenerateToken(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	data.Token = newToken

	if err := u.smtpService.SendVerificationEmail(ctx, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

func (u *UseCases) SendNewRequestMessage(ctx context.Context, code, email string) error {
	data := &dto.EmailData{
		Email:        email,
		Subject:      "Subject: Tu solicitud de aviso de obra ya está en validación - Solicitudes online San Isidro\r\n",
		BodyTemplate: mime,
	}

	if err := u.smtpService.SendNewRequestEmail(ctx, code, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	return nil
}

func (u *UseCases) SendUpdateRequestByCodeMessage(ctx context.Context, code, email, obs string) error {
	data := &dto.EmailData{
		Email:        email,
		Subject:      "Subject: Recibiste observaciones en tu solicitud de Aviso de obra - San Isidro\r\n",
		BodyTemplate: mime,
	}

	if err := u.smtpService.SendUpdateRequestByCodeMessage(ctx, code, obs, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	return nil
}

func (u *UseCases) SendUpdateRequestMessage(ctx context.Context, code, email string) error {
	data := &dto.EmailData{
		Email:        email,
		Subject:      "Subject: Tu solicitud de aviso de obra se actualizó correctamente - San Isidro\r\n",
		BodyTemplate: mime,
	}

	if err := u.smtpService.SendUpdateRequestMessage(ctx, code, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	return nil
}

func (u *UseCases) SendValidateRequestMessage(ctx context.Context, code, email string) error {
	data := &dto.EmailData{
		Email:        email,
		Subject:      "Subject: Tu solicitud de Aviso de Obra se aprobó con éxito! - San Isidro\r\n",
		BodyTemplate: mime,
	}

	if err := u.smtpService.SendValidateRequestMessage(ctx, code, data); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}
	return nil
}

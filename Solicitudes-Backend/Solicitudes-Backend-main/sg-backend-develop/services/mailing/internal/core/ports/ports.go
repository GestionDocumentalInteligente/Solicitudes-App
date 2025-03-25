package ports

import (
	"context"

	sdkjwtdefs "github.com/teamcubation/sg-mailing-api/pkg/jwt/v5/defs"

	dto "github.com/teamcubation/sg-mailing-api/internal/core/dto"
	"github.com/teamcubation/sg-mailing-api/internal/core/entities"
)

type SmtpService interface {
	SendVerificationEmail(context.Context, *dto.EmailData) error
	SendNewRequestEmail(context.Context, string, *dto.EmailData) error
	SendUpdateRequestByCodeMessage(context.Context, string, string, *dto.EmailData) error
	SendUpdateRequestMessage(context.Context, string, *dto.EmailData) error
	SendValidateRequestMessage(context.Context, string, *dto.EmailData) error
}

type UseCases interface {
	InitiateEmailVerification(context.Context, string, string) error
	SendNewRequestMessage(context.Context, string, string) error
	SendUpdateRequestByCodeMessage(context.Context, string, string, string) error
	SendUpdateRequestMessage(context.Context, string, string) error
	SendValidateRequestMessage(context.Context, string, string) error
	ActivateAccount(ctx context.Context, token string) error
	ResendActivationEmail(ctx context.Context, token string) error
	ResendActivationEmailExistingUser(ctx context.Context, email string) error
}

type JwtService interface {
	GenerateToken(context.Context, string) (string, error)
	ValidateToken(context.Context, string) (*sdkjwtdefs.TokenClaims, error)
	GetTokenInfo(context.Context, string) (*sdkjwtdefs.TokenClaims, error)
}

type Repository interface {
	UpdateUser(ctx context.Context, email string) error
	GetuserByEmail(ctx context.Context, email string) (*entities.User, error)
}

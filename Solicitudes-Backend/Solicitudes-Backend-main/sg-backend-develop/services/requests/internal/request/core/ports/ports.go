package ports

import (
	"context"

	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

type UseCases interface {
	GetSuggestions(context.Context, string) ([]domain.Suggestion, error)
	CreateRequestByUserID(context.Context, *domain.Request) error
	GetAllRequestsByUserID(context.Context, int64) ([]domain.Request, error)
	GetAllRequestsByCuil(context.Context, string) ([]domain.Request, error)
	CheckAblOwnership(context.Context, string, int) (bool, error)
	CreateRequestByCuil(context.Context, *domain.Request) error
	UpdateRequest(context.Context, *domain.VerifiedRequest) error
	UpdateRequestByFileNumber(context.Context, *domain.Request) error
	ValidateRequest(context.Context, *domain.ValidateRequest) error
	RequestsVerifications(context.Context, string) (*domain.Verification, error)
	AllRequestsVerifications(context.Context) ([]domain.Verification, error)
	AllRequestsValidations(context.Context) ([]domain.Verification, error)
	DocumentsByCode(context.Context, string) (*domain.Request, []domain.Document, string, error)
	ValidationDocumentsByCode(context.Context, string) (*domain.Request, []domain.Document, []domain.Document, error)
	DocumentByID(context.Context, string) (domain.Document, error)
	GetRequestByID(context.Context, int64) (*domain.Request, error)
	GetRequestByExpCode(context.Context, string) (*domain.Request, error)
}

type Repository interface {
	GetSuggestions(context.Context, string, int64, int64) ([]domain.Suggestion, error)
	CreateRequestByUserID(context.Context, *domain.Request) error
	GetAllRequestsByUserID(context.Context, int64) ([]domain.Request, error)
	GetAllRequestsByCuil(context.Context, string) ([]domain.Request, error)
	CheckAblOwnership(context.Context, string, int) (bool, error)
	CreateRequestByCuil(context.Context, *domain.Request) error
	UpdateRequestWithObservations(context.Context, *domain.VerifiedRequest) (string, string, string, error)
	GetRequestPersonByCuil(context.Context, string) (*domain.Request, error)
	GetRequestByFileNumber(context.Context, string) (*domain.Request, error)
	RequestsVerifications(context.Context, string) (*domain.Verification, error)
	GetAllRequestsVerifications(context.Context) ([]domain.Verification, error)
	GetAllRequestsValidations(context.Context) ([]domain.Verification, error)
	GetDocumentsByCode(context.Context, string) ([]domain.Document, error)
	GetValidationDocumentsByCode(ctx context.Context, id string) ([]domain.Document, error)
	GetDocumentByID(context.Context, string) (domain.Document, error)
	GetRequestByID(ctx context.Context, reqID int64) (*domain.Request, error)
	UpdateRequest(ctx context.Context, id int64, code string) error
	UpdateRequestStatus(ctx context.Context, id int64, status int) error
	UpdateUserRequest(ctx context.Context, req *domain.Request) (string, string, error)
	ValidateRequest(context.Context, *domain.ValidateRequest) (int64, error)
	GetInsuranceDocumentByCode(ctx context.Context, id string, docType int) (string, error)
	GetRequestByExpCode(context.Context, string) (*domain.Request, error)
	GetReplacementIFDocumentsByCode(ctx context.Context, id string, insurance bool) ([]domain.Document, error)
	UpdateVerificationStatus(ctx context.Context, fileNumber string) error
	GetPersonByUserID(ctx context.Context, id int64) (string, error)
	UpdateRequestStatusByFileNumber(ctx context.Context, fileNumber string, status int) error
}

type HttpClient interface {
	SendCreatedRequest(context.Context, *domain.Request) (string, error)
	SendUpdateRequest(ctx context.Context, req *domain.Request) error
	SendVerificationDocument(context.Context, string, string, string, string, string) error
	SendValidationDocument(context.Context, string, string, string) error
	SendEmail(ctx context.Context, code, email string) error
	SendEmailUpdate(ctx context.Context, code, email, observations string) error
	SendEmailUpdateRequest(ctx context.Context, code, email string) error
	SendEmailValidateRequest(ctx context.Context, code, email string) error
}

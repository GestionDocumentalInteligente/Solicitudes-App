package request

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/ports"
)

type useCases struct {
	repository ports.Repository
	httpClient ports.HttpClient
}

func NewUseCases(repository ports.Repository, httpClient ports.HttpClient) ports.UseCases {
	return &useCases{
		repository: repository,
		httpClient: httpClient,
	}
}

func (u *useCases) GetAllRequestsByUserID(ctx context.Context, userID int64) ([]domain.Request, error) {
	reqs, err := u.repository.GetAllRequestsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (u *useCases) GetAllRequestsByCuil(ctx context.Context, cuil string) ([]domain.Request, error) {
	reqs, err := u.repository.GetAllRequestsByCuil(ctx, cuil)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (u *useCases) CreateRequestByUserID(ctx context.Context, req *domain.Request) error {
	if err := u.repository.CreateRequestByUserID(ctx, req); err != nil {
		return err
	}

	return nil
}

func (u *useCases) GetSuggestions(ctx context.Context, inputText string) ([]domain.Suggestion, error) {
	partialSuggestion := parseInput(inputText)

	if partialSuggestion.AddrStreet == "" && partialSuggestion.AddrNum == 0 && partialSuggestion.AblNum == 0 {
		return nil, fmt.Errorf("invalid input: all parameters are empty or zero")
	}

	suggestions, err := u.repository.GetSuggestions(ctx, partialSuggestion.AddrStreet, partialSuggestion.AddrNum, partialSuggestion.AblNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions from repository: %w", err)
	}

	return suggestions, nil
}

func (u *useCases) CheckAblOwnership(ctx context.Context, cuil string, ablNum int) (bool, error) {
	userIsOwner, err := u.repository.CheckAblOwnership(ctx, cuil, ablNum)
	if err != nil {
		return false, fmt.Errorf("error checking ABL ownership: %w", err)
	}

	return userIsOwner, nil
}

func (u *useCases) RequestsVerifications(ctx context.Context, cuil string) (*domain.Verification, error) {
	verification, err := u.repository.RequestsVerifications(ctx, cuil)
	if err != nil {
		return nil, fmt.Errorf("error checking ABL ownership: %w", err)
	}

	return verification, nil
}

func (u *useCases) AllRequestsVerifications(ctx context.Context) ([]domain.Verification, error) {
	verification, err := u.repository.GetAllRequestsVerifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("error checking ABL ownership: %w", err)
	}

	return verification, nil
}

func (u *useCases) AllRequestsValidations(ctx context.Context) ([]domain.Verification, error) {
	verification, err := u.repository.GetAllRequestsValidations(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all requests: %w", err)
	}

	return verification, nil
}

func (u *useCases) DocumentsByCode(ctx context.Context, id string) (*domain.Request, []domain.Document, string, error) {
	request, err := u.repository.GetRequestByFileNumber(ctx, id)
	if err != nil {
		return request, nil, "", fmt.Errorf("error getting request: %w", err)
	}

	docType := 13
	if request.Insurance {
		docType = 15
	}

	gedoCode, err := u.repository.GetInsuranceDocumentByCode(ctx, id, docType)
	if err != nil {
		return request, nil, gedoCode, fmt.Errorf("error getting documents: %w", err)
	}

	documents, err := u.repository.GetDocumentsByCode(ctx, id)
	if err != nil {
		return request, documents, gedoCode, fmt.Errorf("error getting documents: %w", err)
	}

	return request, documents, gedoCode, nil
}

func (u *useCases) DocumentByID(ctx context.Context, id string) (domain.Document, error) {
	verification, err := u.repository.GetDocumentByID(ctx, id)
	if err != nil {
		return verification, fmt.Errorf("error checking document. ID: %s: %w", id, err)
	}

	return verification, nil
}

func (u *useCases) CreateRequestByCuil(ctx context.Context, req *domain.Request) error {
	user, err := u.repository.GetRequestPersonByCuil(ctx, req.Cuil)
	if err != nil {
		return err
	}

	req.Cuil = user.Cuil
	req.Dni = user.Dni
	req.FirstName = user.FirstName
	req.LastName = user.LastName
	req.Email = user.Email
	req.Phone = user.Phone

	err = u.repository.CreateRequestByCuil(ctx, req)
	if err != nil {
		return err
	}

	ctx = context.Background()
	go func(ctx context.Context, req *domain.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic in goroutine: %v", r)
			}
		}()

		code, err := u.httpClient.SendCreatedRequest(ctx, req)
		if err != nil {
			u.repository.UpdateRequestStatus(ctx, req.ID, 2)
		} else if err := u.repository.UpdateRequest(ctx, req.ID, code); err != nil {
			fmt.Println(err.Error())
		} else {
			err := u.httpClient.SendEmail(ctx, code, req.Email)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}(ctx, req)

	return nil
}

func (u *useCases) UpdateRequest(ctx context.Context, req *domain.VerifiedRequest) error {
	user, err := u.repository.GetRequestPersonByCuil(ctx, req.Cuil)
	if err != nil {
		return err
	}

	obsTasks, obsProperty, email, err := u.repository.UpdateRequestWithObservations(ctx, req)
	if err != nil {
		return err
	}

	username := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

	ctx = context.Background()
	go func(ctx context.Context, base64Doc, fileNumber, name, email string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic in goroutine: %v", r)
			}
		}()

		err := u.httpClient.SendVerificationDocument(ctx, base64Doc, fileNumber, req.Reference, req.VerificationType, name)
		if err != nil {
			fmt.Println(err.Error())
		}

		if req.Observations != "" {
			observations := strings.TrimSpace(obsTasks) + "</br>" + strings.TrimSpace(obsProperty)
			err = u.httpClient.SendEmailUpdate(ctx, fileNumber, email, observations)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			if err := u.repository.UpdateRequestStatusByFileNumber(ctx, fileNumber, 1); err != nil {
				fmt.Println(err.Error())
			}
		}
	}(ctx, req.FinalVerificationDocument, req.FileNumber, username, email)

	return nil
}

func (u *useCases) ValidateRequest(ctx context.Context, req *domain.ValidateRequest) error {
	userID, err := u.repository.ValidateRequest(ctx, req)
	if err != nil {
		return err
	}

	if req.IsValid {
		user, err := u.repository.GetRequestPersonByCuil(ctx, req.Cuil)
		if err != nil {
			return err
		}
		username := fmt.Sprintf("%s %s", user.FirstName, user.LastName)

		ctx = context.Background()
		go func(ctx context.Context, content, name string) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("panic in goroutine: %v", r)
				}
			}()

			err := u.httpClient.SendValidationDocument(ctx, req.FileNumber, content, name)
			if err != nil {
				fmt.Println(err.Error())
			}

			email, err := u.repository.GetPersonByUserID(ctx, userID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = u.httpClient.SendEmailValidateRequest(ctx, req.FileNumber, email)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(ctx, req.FileContent, username)
	}

	return nil
}

func (u *useCases) GetRequestByID(ctx context.Context, reqID int64) (*domain.Request, error) {
	req, err := u.repository.GetRequestByID(ctx, reqID)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *useCases) GetRequestByExpCode(ctx context.Context, code string) (*domain.Request, error) {
	req, err := u.repository.GetRequestByExpCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *useCases) ValidationDocumentsByCode(ctx context.Context, id string) (*domain.Request, []domain.Document, []domain.Document, error) {
	request, err := u.repository.GetRequestByFileNumber(ctx, id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting request: %w", err)
	}

	documents, err := u.repository.GetValidationDocumentsByCode(ctx, id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting documents: %w", err)
	}

	replacementIfs, err := u.repository.GetReplacementIFDocumentsByCode(ctx, id, request.Insurance)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting if documents: %w", err)
	}

	for i, doc := range documents {
		if doc.Title == "" {
			if i == 0 {
				documents[i].Title = "Potestad sobre el inmueble"
			} else {
				documents[i].Title = "Tareas a realizar"
			}
		}

		if i == 0 || i == 1 {
			adjustedTime := time.Time(request.VerifyDate).Add(-3 * time.Hour)
			documents[i].VerifiedBy = request.VerifyBy
			documents[i].VerifiedDate = adjustedTime.Format("2006-01-02 15:04")
		} else {
			adjustedTime := time.Time(request.VerifyDateTask).Add(-3 * time.Hour)
			documents[i].VerifiedBy = request.VerifyByTasks
			documents[i].VerifiedDate = adjustedTime.Format("2006-01-02 15:04")
		}
	}

	return request, documents, orderReplacementIFs(replacementIfs), nil
}

func orderReplacementIFs(documents []domain.Document) []domain.Document {
	order := []int{4, 1, 2, 3, 9, 10, 11, 14, 13, 15, 12, 6, 17, 16}

	orderMap := make(map[int]int)
	for i, v := range order {
		orderMap[v] = i
	}

	validDocuments := []domain.Document{}
	for _, doc := range documents {
		if _, exists := orderMap[doc.Type]; exists {
			validDocuments = append(validDocuments, doc)
		}
	}

	sort.SliceStable(validDocuments, func(i, j int) bool {
		return orderMap[validDocuments[i].Type] < orderMap[validDocuments[j].Type]
	})

	return validDocuments
}

func (u *useCases) UpdateRequestByFileNumber(ctx context.Context, req *domain.Request) error {
	user, err := u.repository.GetRequestPersonByCuil(ctx, req.Cuil)
	if err != nil {
		return err
	}

	req.Cuil = user.Cuil
	req.Dni = user.Dni
	req.FirstName = user.FirstName
	req.LastName = user.LastName
	req.Email = user.Email
	req.Phone = user.Phone

	req.Observations, req.ObservationsTasks, err = u.repository.UpdateUserRequest(ctx, req)
	if err != nil {
		return err
	}

	ctx = context.Background()
	go func(ctx context.Context, req *domain.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic in goroutine: %v", r)
			}
		}()

		err := u.httpClient.SendUpdateRequest(ctx, req)
		if err != nil {
			if err := u.repository.UpdateRequestStatus(ctx, req.ID, 7); err != nil {
				fmt.Println(err.Error())
				return
			}
		} else {
			err := u.repository.UpdateVerificationStatus(ctx, req.FileNumber)
			if err != nil {
				fmt.Println(err.Error())
			}

			err = u.httpClient.SendEmailUpdateRequest(ctx, req.FileNumber, req.Email)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}(ctx, req)

	return nil
}

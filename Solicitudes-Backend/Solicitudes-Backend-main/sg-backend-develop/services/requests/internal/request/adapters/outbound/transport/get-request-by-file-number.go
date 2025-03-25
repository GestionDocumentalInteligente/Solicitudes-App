package transport

import (
	"database/sql"

	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

// ToDocumentRequest converts a DocumentModel to a DocumentRequest
func ToDocumentRequest(model *DocumentModel) *domain.DocumentRequest {
	if model == nil {
		return nil
	}

	return &domain.DocumentRequest{
		Name:    model.Filename.String, // Using String from sql.NullString
		Type:    domain.DocumentTypeID(model.Type),
		Content: model.Content,
	}
}

// ToDocumentRequestList converts a slice of DocumentModel to a slice of DocumentRequest
func DocumentModelToDocumentRequestList(models []DocumentModel) []domain.DocumentRequest {
	if len(models) == 0 {
		return []domain.DocumentRequest{}
	}

	requests := make([]domain.DocumentRequest, len(models))
	for i, model := range models {
		request := ToDocumentRequest(&model)
		if request != nil {
			requests[i] = *request
		}
	}

	return requests
}

// Optional: mapper in the opposite direction if needed
func ToDocumentModel(request *domain.DocumentRequest) *DocumentModel {
	if request == nil {
		return nil
	}

	return &DocumentModel{
		Description: sql.NullString{
			String: request.Name,
			Valid:  request.Name != "",
		},
		Type:    int(request.Type),
		Content: request.Content,
	}
}

// Optional: mapper for slice in opposite direction if needed
func DocumentRequestListToDocumentModelList(requests []domain.DocumentRequest) []DocumentModel {
	if len(requests) == 0 {
		return []DocumentModel{}
	}

	models := make([]DocumentModel, len(requests))
	for i, request := range requests {
		model := ToDocumentModel(&request)
		if model != nil {
			models[i] = *model
		}
	}

	return models
}

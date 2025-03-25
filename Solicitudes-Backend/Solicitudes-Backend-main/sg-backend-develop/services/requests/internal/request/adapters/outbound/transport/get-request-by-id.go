package transport

import "github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"

func ToRequestDomain(model *RequestDataModel) *domain.Request {
	if model == nil {
		return nil
	}

	return &domain.Request{
		UserID:             model.UserID,
		UserType:           domain.UserType(model.UserType.String),
		PropertyID:         model.PropertyID,
		ABLDebt:            model.ABLDebt,
		SelectedActivities: model.SelectedActivities,
		EstimatedTime:      model.EstimatedTime,
		Insurance:          model.Insurance,
		FileNumber:         model.FileNumber,
		Documents:          toDocumentRequestDomain(model.Documents),
		StatusName:         model.StatusName,
		CreatedAt:          model.CreatedAt,
		ProjectDesc:        model.Description,
	}
}

// Helper function para convertir los documentos
func toDocumentRequestDomain(docs []DocumentDataModel) []domain.DocumentRequest {
	if docs == nil {
		return nil
	}

	domainDocs := make([]domain.DocumentRequest, len(docs))
	for i, doc := range docs {
		domainDocs[i] = domain.DocumentRequest{
			Name:    doc.Name,
			Type:    domain.DocumentTypeID(doc.Type),
			Content: doc.Content,
		}
	}
	return domainDocs
}

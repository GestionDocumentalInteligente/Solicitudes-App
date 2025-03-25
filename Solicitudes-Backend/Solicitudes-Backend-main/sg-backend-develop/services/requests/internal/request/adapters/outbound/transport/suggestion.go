package transport

import (
	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

type SuggestionDataModel struct {
	ID         int64  `db:"property_id"`
	AddrStreet string `db:"street"`
	AddrNum    int64  `db:"number"`
	AblNum     int64  `db:"abl_number"`
}

// ToSuggestionDomainList convierte una lista de Suggestion (modelo de datos) a una lista de domain.Suggestion (entidad de dominio)
func ToSuggestionDomainList(dataModels []SuggestionDataModel) []domain.Suggestion {
	var suggestions []domain.Suggestion

	for _, model := range dataModels {
		suggestions = append(suggestions, domain.Suggestion{
			ID:         model.ID,
			AddrStreet: model.AddrStreet,
			AddrNum:    model.AddrNum,
			AblNum:     model.AblNum,
		})
	}

	return suggestions
}

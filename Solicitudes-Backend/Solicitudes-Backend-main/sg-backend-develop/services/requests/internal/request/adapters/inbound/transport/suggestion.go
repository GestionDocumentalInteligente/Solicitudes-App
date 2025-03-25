// internal/core/request/adapters/inbound/transport/suggestion.go
package transport

import (
	"github.com/teamcubation/sg-backend/services/requests/internal/request/core/domain"
)

type SuggestionJson struct {
	ID         int64  `json:"property_id"`
	AddrStreet string `json:"address_street"`
	AddrNum    int64  `json:"address_number"`
	AblNum     int64  `json:"abl_number"`
}

// presenter
func ToSuggestionPresenter(list []domain.Suggestion) []SuggestionJson {
	suggestions := make([]SuggestionJson, len(list))
	for i, model := range list {
		suggestions[i] = SuggestionJson{
			ID:         model.ID,
			AddrStreet: model.AddrStreet,
			AddrNum:    model.AddrNum,
			AblNum:     model.AblNum,
		}
	}
	return suggestions
}

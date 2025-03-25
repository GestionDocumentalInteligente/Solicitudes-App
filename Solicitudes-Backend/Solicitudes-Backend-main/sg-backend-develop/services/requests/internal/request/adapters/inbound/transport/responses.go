// internal/core/request/adapters/inbound/transport/responses.go
package transport

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type SuggestionsResponse struct {
	Suggestions []SuggestionJson `json:"suggestions"`
}

type AblOwnershipResponse struct {
	AblOwnership bool `json:"abl_ownership"`
}

type VerificationResponse struct {
	Verification *VerificationRequestPresenter `json:"verification_request"`
}

type RequestsResponse struct {
	Requests []GetAllReqPresenterRequest `json:"requests"`
}

type DocumentsResponse struct {
	Documents []Document `json:"documents"`
}

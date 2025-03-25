// internal/core/request/adapters/inbound/transport/errors.go
package transport

import "errors"

var (
	ErrInvalidPayload     = errors.New("invalid request payload")
	ErrMissingQueryParam  = errors.New("missing required query parameter")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidABLNumber   = errors.New("invalid ABL number")
	ErrInternalServer     = errors.New("internal server error")
	ErrNoSuggestionsFound = errors.New("no suggestions found")
)

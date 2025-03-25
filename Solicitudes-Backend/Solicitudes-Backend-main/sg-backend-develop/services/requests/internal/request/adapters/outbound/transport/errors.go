package transport

import "errors"

var (
	ErrNilRequest      = errors.New("request cannot be nil")
	ErrEmptyAddress    = errors.New("address string cannot be empty")
	ErrInvalidNumber   = errors.New("invalid address number")
	ErrAddressNotFound = errors.New("address not found")
	ErrCreateRequest   = errors.New("failed to create request")
	ErrPersonNotFound  = errors.New("person not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrCreateDocument  = errors.New("error creating document")
	ErrUpdateRequest   = errors.New("error updating request")
)

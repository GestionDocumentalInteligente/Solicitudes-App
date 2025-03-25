package entities

import "fmt"

type ErrorTokenInvalid struct {
	Msg string
}

func (e *ErrorTokenInvalid) Error() string {
	return fmt.Sprintf("token invalid: %s", e.Msg)
}

type NotFoundInDatabase struct {
	Msg string
}

func (e *NotFoundInDatabase) Error() string {
	return fmt.Sprintf("database error: %s", e.Msg)
}

type UserAlreadyActive struct {
	Email string
}

func (e *UserAlreadyActive) Error() string {
	return fmt.Sprintf("User already active: %s", e.Email)
}

package file

import "net/http"

type CustomError struct {
	StatusCode int
	Message    string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(statusCode int, message string) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func BadRequest(message string) *CustomError {
	return NewCustomError(http.StatusBadRequest, message)
}

func NotFound(message string) *CustomError {
	return NewCustomError(http.StatusNotFound, message)
}

func InternalServerError(message string) *CustomError {
	return NewCustomError(http.StatusInternalServerError, message)
}

package fixer

import "net/http"

// NewError creates a new Error
func NewError(msg string) *Error {
	return &Error{msg: msg}
}

// Error type for Fixer API requests
type Error struct {
	msg string
}

// Error message
func (e *Error) Error() string {
	return e.msg
}

// Errors
var (
	ErrNotFound            = NewError(http.StatusText(http.StatusNotFound))
	ErrUnprocessableEntity = NewError(http.StatusText(http.StatusUnprocessableEntity))
)

func responseError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnprocessableEntity:
		return ErrUnprocessableEntity
	default:
		return NewError("Unexpected status: " + resp.Status)
	}
}
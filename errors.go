package fixer

import (
	"fmt"
	"io"
	"net/http"
)

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
	ErrNilResponse         = NewError("Unexpected nil response")
	ErrUnexpectedStatus    = NewError("Unexpected status")
	ErrNotFound            = NewError(http.StatusText(http.StatusNotFound))
	ErrUnprocessableEntity = NewError(http.StatusText(http.StatusUnprocessableEntity))
	ErrUnauthorized        = NewError(http.StatusText(http.StatusUnauthorized))
	ErrInternalServerError = NewError(http.StatusText(http.StatusInternalServerError))
)

func responseError(resp *http.Response) error {
	if resp == nil {
		return ErrNilResponse
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusUnprocessableEntity:
		return ErrUnprocessableEntity
	case http.StatusInternalServerError:
		return ErrInternalServerError
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Response Error: %s: %s", resp.Status, string(body))
	}
}

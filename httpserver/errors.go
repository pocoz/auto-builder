package httpserver

import (
	"fmt"
)

// Error is a tms service error.
type Error struct {
	Kind    ErrorKind
	Message string
}

// ErrorResponse HTTP answer when an error occurs
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description, omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// ErrorKind is a kind of the service error.
type ErrorKind uint8

// Error kinds.
const (
	ErrBadParams ErrorKind = iota
	ErrNotFound
	ErrConflict
	ErrInternal
	ErrUnauthorized
	ErrForbidden
	ErrNotAcceptable
	ErrNotAllowed
)

func errorf(kind ErrorKind, format string, v ...interface{}) error {
	return &Error{
		Kind:    kind,
		Message: fmt.Sprintf(format, v...),
	}
}

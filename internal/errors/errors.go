package errors

import (
	"fmt"
)

// CustomError defines a custom error type for the application.
// It can wrap an underlying error to preserve original context.
type CustomError struct {
	Message string
	Err     error
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return fmt.Sprintf("CustomError: %s", e.Message)
	}
	return fmt.Sprintf("CustomError: %s: %v", e.Message, e.Err)
}

// Unwrap returns the underlying error, if any, to allow errors.Is/As checks.
func (e *CustomError) Unwrap() error { return e.Err }

// New creates a CustomError wrapping an optional underlying error.
func New(msg string, err error) error {
	return &CustomError{Message: msg, Err: err}
}

package errors

import "fmt"

// CustomError defines a custom error type for the application.
type CustomError struct {
	Message string
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	return fmt.Sprintf("CustomError: %s", e.Message)
}

package domain

import "fmt"

// ErrorCode is a custom type for defining application-level error categories
type ErrorCode string

const (
	CodeNotFound   ErrorCode = "NOT_FOUND"
	CodeConflict   ErrorCode = "CONFLICT"
	CodeInternal   ErrorCode = "INTERNAL"
	CodeValidation ErrorCode = "VALIDATION"
)

// AppError is the custom domain error
type AppError struct {
	Code    ErrorCode
	Message string
	Field   string
	Err     error // The original error for logging
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows this error to work with errors.As and errors.Is
func (e *AppError) Unwrap() error {
	return e.Err
}

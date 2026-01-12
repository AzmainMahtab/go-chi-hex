// Package domain
// This one deals with app level errors
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

type ErrorItem struct {
	Field   string
	Message string
}

// AppError is the custom domain error
type AppError struct {
	Code    ErrorCode
	Message string
	Errors  []ErrorItem
	Field   string
	Err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows this error to work with errors.As and errors.Is
func (e *AppError) Unwrap() error {
	return e.Err
}

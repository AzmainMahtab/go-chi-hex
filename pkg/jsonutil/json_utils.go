// Package jsonutil
// jsonutil contains all the JSON marshaling and unmarshaling
package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Response is your standard API envelope
type Response struct {
	Success    bool        `json:"success"`
	Status     string      `json:"status"`            // "success", "fail", or "error"
	StatusCode int         `json:"statusCode"`        // HTTP status code
	Message    string      `json:"message,omitempty"` // Human-readable message
	Data       interface{} `json:"data,omitempty"`    // Primary response data
	Meta       interface{} `json:"meta,omitempty"`    // Pagination, etc.
	Errors     []ErrorItem `json:"errors,omitempty"`  // Error details
}

// ErrorItem represents machine-readable error details
type ErrorItem struct {
	Code    string `json:"code,omitempty"`    // e.g., "INVALID_EMAIL"
	Field   string `json:"field,omitempty"`   // e.g., "email"
	Message string `json:"message,omitempty"` // User-friendly description
}

const MaxRequestBodySize = 1048576

// ReadJSON decodes a single JSON object from the request body
func ReadJSON(w http.ResponseWriter, r *http.Request, target any) error {
	// 1. Defer closing the original body.

	defer r.Body.Close()

	// 2. Limit the body size and replace the request body with the limited reader.

	limitedReader := http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	decoder := json.NewDecoder(limitedReader)

	// 3. Enforce strict JSON rules

	decoder.DisallowUnknownFields()

	err := decoder.Decode(target)
	if err != nil {

		// 4. Detailed error handling for production-grade APIs

		var syntaxError *json.SyntaxError

		var unmarshalTypeError *json.UnmarshalTypeError

		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxError):

			return fmt.Errorf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)

		case errors.Is(err, http.ErrNotSupported): // Catches the error from MaxBytesReader if limit is exceeded

			return fmt.Errorf("Request body must not be larger than %d bytes", MaxRequestBodySize)

		case errors.As(err, &unmarshalTypeError):

			// Handle wrong data type for a field

			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Request body contains an incorrect type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("Request body contains an incorrect type (at position %d)", unmarshalTypeError.Offset)

		case errors.As(err, &invalidUnmarshalError):

			// This is usually a programming error (passing nil or non-pointer to target)

			panic(err)

		case strings.Contains(err.Error(), "json: unknown field"):

			// Handle disallowed unknown fields

			fieldName := strings.Split(err.Error(), "\"")[1]

			return fmt.Errorf("Request body contains unknown field %s", fieldName)

		default:

			return err

		}

	}

	// 5. Check if there are any remaining JSON elements after the first object

	if decoder.More() {
		return errors.New("Request body must contain only a single JSON object")
	}

	return nil
}

// WriteJSON handles all successful responses
func WriteJSON(w http.ResponseWriter, status int, data any, meta any, message string) error {
	resp := Response{
		Success:    true,
		Status:     "success",
		StatusCode: status,
		Message:    message,
		Data:       data,
		Meta:       meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(resp)
}

// ErrorResponse is the generic helper for sending errors
func ErrorResponse(w http.ResponseWriter, status int, message string, errs []ErrorItem) {
	// Map HTTP status to "fail" (4xx) or "error" (5xx)
	statusText := "error"
	if status >= 400 && status < 500 {
		statusText = "fail"
	}

	resp := Response{
		Success:    false,
		Status:     statusText,
		StatusCode: status,
		Message:    message,
		Errors:     errs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// --- Specific Error Shortcuts ---

// BadRequestResponse for validation or malformed JSON
func BadRequestResponse(w http.ResponseWriter, message string, errs []ErrorItem) {
	ErrorResponse(w, http.StatusBadRequest, message, errs)
}

// ServerErrorResponse shields internal details.
// You pass the ACTUAL error for logging, but the user sees a generic message.
func ServerErrorResponse(w http.ResponseWriter, loggerErr error) {
	// In a real app, you'd log loggerErr here
	// log.Printf("INTERNAL ERROR: %v", loggerErr)

	message := "An unexpected server error occurred"
	ErrorResponse(w, http.StatusInternalServerError, message, nil)
}

// NotFoundResponse for missing resources
func NotFoundResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message, nil)
}

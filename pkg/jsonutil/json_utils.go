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

// APIResponse is the standard "Envelope" for all successful responses
type APIResponse struct {
	Status  string `json:"status"`            // "success" or "error"
	Data    any    `json:"data,omitempty"`    // The actual payload
	Message string `json:"message,omitempty"` // A human-readable message
}

// APIError defines a structured error for the client
type APIError struct {
	Code    int    `json:"code"`             // Internal application code or HTTP status
	Message string `json:"message"`          // User-friendly error message
	Detail  any    `json:"detail,omitempty"` // Specifics (e.g., validation errors)
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

// WriteJSON handles the consistent response logic
func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Standardize the output into our envelope
	response := APIResponse{
		Status: "success",
		Data:   data,
	}

	js, err := json.Marshal(response)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

// ---  ERROR HELPERS ---

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string, details any) {
	resp := APIResponse{
		Status:  "error",
		Message: message,
		Data: APIError{
			Code:    status,
			Message: message,
			Detail:  details,
		},
	}

	js, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// Log the actual error here for internal tracking
	ErrorResponse(w, r, http.StatusInternalServerError, "An unexpected server error occurred", nil)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error(), nil)
}

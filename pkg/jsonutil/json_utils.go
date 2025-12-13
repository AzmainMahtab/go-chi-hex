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

// Define the maximum allowed size for the request body (1MB)
const MaxRequestBodySize = 1048576

// ReadJSON decodes a single JSON object from the request body into the target struct.
// It enforces limits on body size and disallows unknown fields.
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

// WriteJSON encodes data into JSON, sets the Content-Type, and sends the response.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	// 1. Marshal the data
	js, err := json.Marshal(data)
	if err != nil {
		// Log the internal encoding error (e.g., if data cannot be marshaled)
		// and return an internal server error to the client.
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 2. Set headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// 3. Write the JSON bytes
	if _, err = w.Write(js); err != nil {
		// Log if there was an issue writing the response bytes (e.g., connection prematurely closed)
		return fmt.Errorf("failed to write JSON response: %w", err)
	}

	return nil
}

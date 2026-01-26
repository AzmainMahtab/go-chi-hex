// Package handlers
// THis one contains the helpers for handler
package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ReadIDParam Read the id parameter from the userl
func ReadIDParam(r *http.Request) (string, error) {
	idStr := chi.URLParam(r, "id")
	// id, err := strconv.Atoi(idStr)
	//
	// if err != nil || id < 1 {
	// 	return 0, &domain.AppError{
	// 		Code:    domain.CodeValidation,
	// 		Message: "invalid resource identifier",
	// 	}
	// }
	return idStr, nil
}

// ParseQueryInt safely converts a URL query parameter to an integer with a default fallback
func ParseQueryInt(r *http.Request, key string, defaultValue int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}

	return val
}

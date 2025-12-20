// Package handlers
// THis one contains the helpers for handler
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ReadIDPeram(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

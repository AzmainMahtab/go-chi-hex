// Package handlers
// THis one contains the helpers for handler
package handlers

import (
	"net/http"
	"strconv"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/go-chi/chi/v5"
)

func ReadIDParam(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id < 1 {
		return 0, &domain.AppError{
			Code:    domain.CodeValidation,
			Message: "invalid resource identifier",
		}
	}
	return id, nil
}

package handlers

import (
	"errors"
	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/pkg/jsonutil"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError

	if errors.As(err, &appErr) {
		switch appErr.Code {
		case domain.CodeConflict:
			// Convert domain items to jsonutil items
			items := make([]jsonutil.ErrorItem, 0)
			for _, e := range appErr.Errors {
				items = append(items, jsonutil.ErrorItem{Field: e.Field, Message: e.Message})
			}
			// If slice was empty (single field from DB), use the single Field
			if len(items) == 0 && appErr.Field != "" {
				items = append(items, jsonutil.ErrorItem{Field: appErr.Field, Message: appErr.Message})
			}
			jsonutil.ConflictResponse(w, appErr.Message, items)

		case domain.CodeNotFound:
			jsonutil.NotFoundResponse(w, appErr.Message)

		case domain.CodeValidation:
			jsonutil.BadRequestResponse(w, appErr.Message, nil)

		default:
			jsonutil.ServerErrorResponse(w, appErr.Err)
		}
		return
	}

	jsonutil.ServerErrorResponse(w, err)
}

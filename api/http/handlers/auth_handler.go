// Package handlers
// this one contains the auth handlers
package handlers

import (
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/apiutil"
	"github.com/AzmainMahtab/go-chi-hex/api/http/dto"
	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/AzmainMahtab/go-chi-hex/pkg/jsonutil"
)

type AuthHandler struct {
	svc ports.AuthService
}

func NewAuthHandler(svc ports.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest

	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	if errs := apiutil.ValidateStruct(req); errs != nil {
		jsonutil.BadRequestResponse(w, "Invalid data", errs)
		return
	}

	authLogin := &domain.AuthLogin{
		Email:    req.Email,
		Password: req.Password,
	}

	token, err := a.svc.Login(r.Context(), *authLogin)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, token, nil, "Login success")

}

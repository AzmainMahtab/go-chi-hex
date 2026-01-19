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

// Login handles user authentication and returns a JWT token.
// @Summary      User Login
// @Description  Authenticate user with email and password to receive a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AuthRequest  true  "Login Credentials"
// @Success      200      {object}  map[string]interface{} "Login success"
// @Failure      400      {object}  map[string]interface{} "Bad request or invalid data"
// @Failure      401      {object}  map[string]interface{} "Unauthorized"
// @Failure      500      {object}  map[string]interface{} "Internal server error"
// @Router       /auth/login [post]
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

// Package handlers
// this one contains the auth handlers
package handlers

import (
	"log"
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/apiutil"
	"github.com/AzmainMahtab/go-chi-hex/api/http/dto"
	"github.com/AzmainMahtab/go-chi-hex/api/http/middleware"
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

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterUserRequest  true  "User Data"
// @Success      201   {object}  dto.UserResponse
// @Router       /auth/register [post]
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.BadRequestResponse(w, "Bad Request", nil)
		return
	}

	if errs := apiutil.ValidateStruct(req); errs != nil {
		jsonutil.BadRequestResponse(w, "Invalid data", errs)
		return
	}

	// Convert DTO to Domain for Service Input
	userDomain := domain.User{
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	res, err := a.svc.Register(r.Context(), userDomain)
	if err != nil {
		HandleError(w, err)
		log.Printf("SVC ERR: %v", err)
		return
	}

	// Map Domain back to DTO for Response
	jsonutil.WriteJSON(w, http.StatusCreated, a.mapToResponse(res), nil, "User registered successfully")
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

// Logout revokes the refresh token.
// @Summary      Logout User
// @Description  Blacklists the provided refresh token to end the session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.LogoutRequest  true  "Refresh Token to revoke"
// @Success      200      {object}  jsonutil.Response "Logout success"
// @Failure      401      {object}  jsonutil.Response "Unauthorized"
// @Router       /auth/logout [post]
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest

	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(domain.UserClaims)
	if !ok {
		jsonutil.UnauthorizedResponse(w, "Unauthorized: No claims found")
		return
	}

	if err := a.svc.Logout(r.Context(), req.RefreshToken, claims); err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, nil, nil, "Logout success")
}

// Rotate provides new tokens using a refresh token.
// @Summary      Rotate Tokens
// @Description  Generates a new access and refresh token pair using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RotateRequest  true  "Refresh Token"
// @Success      200      {object}  domain.Tokenpair "Rotation success"
// @Failure      401      {object}  jsonutil.Response "Token revoked or invalid"
// @Router       /auth/rotate [post]
func (a *AuthHandler) Rotate(w http.ResponseWriter, r *http.Request) {
	var req dto.RotateRequest

	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	tokenPair, err := a.svc.Rotate(r.Context(), req.RefreshToken)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, tokenPair, nil, "Tokens rotated successfully")
}

func (a *AuthHandler) mapToResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         u.UUID,
		UserName:   u.UserName,
		Email:      u.Email,
		Phone:      u.Phone,
		UserStatus: u.UserStatus,
		CreatedAt:  u.CreatedAt,
	}
}

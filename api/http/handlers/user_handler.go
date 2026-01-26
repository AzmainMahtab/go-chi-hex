// Package handlers
// This one holds all the handler relate to user
package handlers

import (
	"log"
	"net/http"

	"github.com/AzmainMahtab/go-chi-hex/api/http/apiutil"
	"github.com/AzmainMahtab/go-chi-hex/api/http/dto"
	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/AzmainMahtab/go-chi-hex/internal/ports"
	"github.com/AzmainMahtab/go-chi-hex/pkg/jsonutil"
)

type UserHandler struct {
	svc ports.UserService
}

func NewUserHandler(svc ports.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register godoc
// @Summary      Register a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterUserRequest  true  "User Data"
// @Success      201   {object}  dto.UserResponse
// @Router       /user [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.svc.RegisterUser(r.Context(), userDomain)
	if err != nil {
		HandleError(w, err)
		log.Printf("SVC ERR: %v", err)
		return
	}

	// Map Domain back to DTO for Response
	jsonutil.WriteJSON(w, http.StatusCreated, h.mapToResponse(res), nil, "User registered successfully")
}

// List godoc
// @Summary      List active users
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  dto.UserResponse
// @Router       /user [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	// Extract and convert query parameters directly into the Domain Filter
	filter := domain.UserFilter{
		UserName:    r.URL.Query().Get("user_name"),
		Email:       r.URL.Query().Get("email"),
		Phone:       r.URL.Query().Get("phone"),
		UserStatus:  r.URL.Query().Get("user_status"),
		ShowDeleted: false,                         // Explicitly false for the active list
		Limit:       ParseQueryInt(r, "limit", 10), // Default to 10
		Offset:      ParseQueryInt(r, "offset", 0), // Default to 0
	}

	// Call Service
	users, err := h.svc.ListUsers(r.Context(), filter)
	if err != nil {
		HandleError(w, err)
		return
	}

	//  Respond with the typed domain slice
	jsonutil.WriteJSON(w, http.StatusOK, h.mapSliceToResponse(users), nil, "Active users retrieved")
}

// GetByID godoc
// @Summary      Get user by ID
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Router       /user/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDParam(r)
	if err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, h.mapToResponse(user), nil, "User fetched")
}

// Update godoc
// @Summary      Update user partially
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path      string                      true  "User ID"
// @Param        user  body      dto.UpdateUserRequest true  "Fields to update"
// @Security     BearerAuth
// @Success      200   {object}  dto.UserResponse
// @Router       /user/{id} [patch]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDParam(r)
	if err != nil {
		HandleError(w, err)
		return
	}

	// Decode JSON into DTO
	var req dto.UpdateUserRequest
	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	if errs := apiutil.ValidateStruct(req); errs != nil {
		jsonutil.BadRequestResponse(w, "Invalid data", errs)
		return
	}

	// Map DTO to Domain.UserUpdate (Strictly Typed)
	updateParams := domain.UserUpdate{
		ID:       id,
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}

	// Execute Service
	updatedUser, err := h.svc.UpdateUser(r.Context(), updateParams)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, h.mapToResponse(updatedUser), nil, "User updated successfully")
}

// Remove godoc
// @Summary      Soft delete user
// @Tags         user
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      204  "No Content"
// @Router       /user/{id} [delete]
func (h *UserHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDParam(r)
	if err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	if err := h.svc.RemoveUser(r.Context(), id); err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusNoContent, nil, nil, "User moved to trash")
}

// Restore godoc
// @Summary      Restore a delete user
// @Description  Restore a user that has been soft deleted
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  string "Invalid ID"
// @Failure      500  {object}  string "Internal Server Error"
// @Router       /user/{id}/restore [patch]
func (h *UserHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDParam(r)
	if err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	user, err := h.svc.RestoreUser(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, h.mapToResponse(user), nil, "User restored"); err != nil {
		log.Printf("Handler: GetTrashed error: %v", err)
		HandleError(w, err)
		return
	}
}

// GetTrashed godoc
// @Summary      List soft-deleted users
// @Description  Retrieves all users where deleted_at is not null
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.UserResponse
// @Router       /user/trash [get]
func (h *UserHandler) GetTrashed(w http.ResponseWriter, r *http.Request) {
	filter := domain.UserFilter{
		UserName:    r.URL.Query().Get("user_name"),
		Email:       r.URL.Query().Get("email"),
		ShowDeleted: true, // Internal logic for the repository
		Limit:       ParseQueryInt(r, "limit", 10),
		Offset:      ParseQueryInt(r, "offset", 0),
	}

	//  Call the dedicated Trash service method
	users, err := h.svc.GetTrashedUsers(r.Context(), filter)
	if err != nil {
		HandleError(w, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, h.mapSliceToResponse(users), nil, "Secret trash retrieved")
}

// Prune godoc
// @Summary      Permanently delete a user
// @Description  Hard deletes a user record from the database. This action cannot be undone.
// @Tags         user
// @Security 		 BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      204  {string}  string "User permanently deleted"
// @Failure      400  {object}  string "Invalid ID"
// @Failure      500  {object}  string "Internal Server Error"
// @Router       /user/{id}/prune [delete]
func (h *UserHandler) Prune(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDParam(r)
	if err != nil {
		jsonutil.BadRequestResponse(w, "Bad request", nil)
		return
	}

	if err := h.svc.PermanentlyDeleteUser(r.Context(), id); err != nil {
		HandleError(w, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusNoContent, nil, nil, "User permanently deleted"); err != nil {
		HandleError(w, err)
		return
	}
}

// --- MAPPING HELPERS ---

func (h *UserHandler) mapToResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         u.UUID,
		UserName:   u.UserName,
		Email:      u.Email,
		Phone:      u.Phone,
		UserStatus: u.UserStatus,
		CreatedAt:  u.CreatedAt,
	}
}

func (h *UserHandler) mapSliceToResponse(users []*domain.User) []dto.UserResponse {
	res := make([]dto.UserResponse, len(users))
	for i, u := range users {
		res[i] = h.mapToResponse(u)
	}
	return res
}

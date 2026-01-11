// Package handlers
// This one holds all the handler relate to user
package handlers

import (
	"log"
	"net/http"

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
// @Success      200  {array}  dto.UserResponse
// @Router       /user [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context(), nil)
	if err != nil {
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, h.mapSliceToResponse(users), nil, "User list fetched successfully")
}

// GetByID godoc
// @Summary      Get user by ID
// @Tags         user
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Router       /user/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDPeram(r)
	if err != nil {
		return
	}

	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, h.mapToResponse(user), nil, "User fetched")
}

// Update godoc
// @Summary      Update user partially
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "User ID"
// @Param        user  body      dto.UpdateUserRequest true  "Fields to update"
// @Success      200   {object}  dto.UserResponse
// @Router       /user/{id} [patch]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDPeram(r)
	if err != nil {
		return
	}

	var updates map[string]any
	if err := jsonutil.ReadJSON(w, r, &updates); err != nil {
		return
	}

	// IMPORTANT: remove keys that shouldn't be updated via API
	delete(updates, "id")
	delete(updates, "created_at")

	user, err := h.svc.UpdateUser(r.Context(), id, updates)
	if err != nil {
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, h.mapToResponse(user), nil, "User updated successfully")
}

// Remove godoc
// @Summary      Soft delete user
// @Tags         user
// @Param        id   path      int  true  "User ID"
// @Success      204  "No Content"
// @Router       /user/{id} [delete]
func (h *UserHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDPeram(r)
	if err != nil {
		return
	}

	if err := h.svc.RemoveUser(r.Context(), id); err != nil {
		return
	}

	jsonutil.WriteJSON(w, http.StatusNoContent, nil, nil, "User moved to trash")
}

// Restore godoc
// @Summary      Restore a delete user
// @Description  Restore a user that has been soft deleted
// @Tags         user
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  string "Invalid ID"
// @Failure      500  {object}  string "Internal Server Error"
// @Router       /user/{id}/restore [patch]
func (h *UserHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDPeram(r)
	if err != nil {
		return
	}

	user, err := h.svc.RestoreUser(r.Context(), id)
	if err != nil {
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, h.mapToResponse(user), nil, "User restored"); err != nil {
		log.Printf("Handler: GetTrashed error: %v", err)
		return
	}
}

// GetTrashed godoc
// @Summary      List soft-deleted users
// @Description  Retrieves all users where deleted_at is not null
// @Tags         user
// @Produce      json
// @Success      200  {array}   dto.UserResponse
// @Router       /user/trash [get]
func (h *UserHandler) GetTrashed(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetTrashedUsers(r.Context())
	if err != nil {
		log.Printf("Handler: GetTrashed error: %v", err)
		//  jsonutil handles the error response
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, h.mapSliceToResponse(users), nil, "Trash fetched successfully"); err != nil {
		log.Printf("Handler: WriteJSON error: %v", err)
		return
	}
}

// Prune godoc
// @Summary      Permanently delete a user
// @Description  Hard deletes a user record from the database. This action cannot be undone.
// @Tags         user
// @Param        id   path      int  true  "User ID"
// @Success      204  {string}  string "User permanently deleted"
// @Failure      400  {object}  string "Invalid ID"
// @Failure      500  {object}  string "Internal Server Error"
// @Router       /user/{id}/prune [delete]
func (h *UserHandler) Prune(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDPeram(r)
	if err != nil {
		return
	}

	if err := h.svc.PermanentlyDeleteUser(r.Context(), id); err != nil {
		log.Printf("Handler: Prune error for ID %d: %v", id, err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusNoContent, nil, nil, "User permanently deleted"); err != nil {
		return
	}
}

// --- MAPPING HELPERS ---

func (h *UserHandler) mapToResponse(u *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         u.ID,
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

// Package handlers
// This one holds all the handler relate to user
package handlers

import (
	"log"
	"net/http"

	"github.com/AzmainMahtab/docpad/api/http/dto"
	"github.com/AzmainMahtab/docpad/internal/ports"
	"github.com/AzmainMahtab/docpad/pkg/jsonutil"
)

type UserHandler struct {
	svc ports.UserService
}

func NewUSerHandler(svc ports.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register godoc
// @Summary      Register a new user !
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      dto.RegisterUserRequest  true  "User Data"
// @Success      201   {object}  dto.UserResponse
// @Router       /user [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest

	if err := jsonutil.ReadJSON(w, r, &req); err != nil {
		return
	}

	res, err := h.svc.RegisterUser(r.Context(), req)
	if err != nil {
		log.Printf("SVC ERR: %v", err)
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusCreated, res, nil); err != nil {
		return
	}
}

// List godoc
// @Summary      List all users (Test!)
// @Description  Returns a test status message for the user service
// @Tags         user
// @Produce      json
// @Success      200  {object} dto.UserResponse
// @Router       /user [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := ReadIDPeram(r)
	if err != nil {
		return
	}

	res, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		return
	}

	if err := jsonutil.WriteJSON(w, http.StatusOK, res, nil); err != nil {
		return
	}
}

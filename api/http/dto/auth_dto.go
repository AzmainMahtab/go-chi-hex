// Package dto
// this one has the auth dto
package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Package dto
// this one has the auth dto
package dto

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email" example:"hehe@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"hehe1234"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJFUzI1NiIsInR5c..."`
}

type RotateRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJFUzI1NiIsInR5c..."`
}

// Package dto
// COntains reques and response shapes for User
package dto

import "time"

// RegisterUserRequest is what the client sends
type RegisterUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,e164"` // e164 ensures international phone format
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest is what the handler recives
type UpdateUserRequest struct {
	UserName *string `json:"user_name,omitempty" validate:"omitempty,min=3,max=32"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
}

// UserResponse is what we send back
type UserResponse struct {
	ID         int       `json:"id"`
	UserName   string    `json:"user_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	UserStatus string    `json:"user_status"`
	CreatedAt  time.Time `json:"created_at"`
}

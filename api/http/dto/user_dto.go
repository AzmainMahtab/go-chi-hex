// Package dto
// COntains reques and response shapes for User
package dto

import "time"

// RegisterUserRequest is what the client sends
type RegisterUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	UserName *string `json:"user_name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Status   *string `json:"status,omitempty"`
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

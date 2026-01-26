// Package dto
// COntains reques and response shapes for User
package dto

import "time"

// RegisterUserRequest is what the client sends
type RegisterUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=32" example:"hehe"`
	Email    string `json:"email" validate:"required,email" example:"hehe@hehemail.com"`
	Phone    string `json:"phone" validate:"required,e164" example:"+8801700000000"` // e164 ensures international phone format
	Password string `json:"password" validate:"required,min=8" example:"Very$tr0ngP@$$w0Rd"`
}

// UpdateUserRequest is what the handler recives
type UpdateUserRequest struct {
	UserName *string `json:"user_name,omitempty" validate:"omitempty,min=3,max=32" example:"hehe"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email" example:"hehe@hehemail.com"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,e164" example:"+8801700000000"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended" example:"active"`
}

// UserResponse is what we send back
type UserResponse struct {
	ID         string    `json:"id"`
	UserName   string    `json:"user_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	UserStatus string    `json:"user_status"`
	CreatedAt  time.Time `json:"created_at"`
}

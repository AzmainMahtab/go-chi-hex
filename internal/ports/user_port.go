// Package ports
// ports contain all the interfaces needed for inbound
// and outbound
package ports

import (
	"context"

	"github.com/AzmainMahtab/docpad/api/http/dto"
	"github.com/AzmainMahtab/docpad/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error

	// Read active users take optional filtering return: list of users and error if any
	ReadAll(ctx context.Context, filter map[string]any, showDeleted bool) ([]*domain.User, error)

	// ReadOne reads a single active user
	ReadOne(ctx context.Context, id int) (*domain.User, error)

	// Update provided fields and update pertially
	Update(ctx context.Context, id int, updates map[string]any) error

	// SoftDelete soft deletes a user set deleted_at current and user_status = 'inactive'
	SoftDelete(ctx context.Context, id int) error

	// Trash lets you read soft deleted users with optional filtering
	Trash(ctx context.Context, filter map[string]any) ([]*domain.User, error)

	// Pruen permanently delets a user
	Prune(ctx context.Context, id int) error
}

type UserService interface {
	RegisterUser(ctx context.Context, req domain.User) (*domain.User, error)

	// ListUsers retrieves users based on filters provided in the request.
	ListUsers(ctx context.Context, filters map[string]any) ([]*domain.User, error)

	// GetUser retrieves a single active user by their unique ID.
	GetUser(ctx context.Context, id int) (*domain.User, error)

	// UpdateUser performs a partial update on a user's information.
	UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*domain.User, error)

	// RemoveUser soft-deletes a user from the active system.
	RemoveUser(ctx context.Context, id int) error

	// GetTrashedUsers retrieves users that have been soft-deleted.
	GetTrashedUsers(ctx context.Context) ([]*domain.User, error)

	// PermanentlyDeleteUser removes a user record from the database entirely.
	PermanentlyDeleteUser(ctx context.Context, id int) error
}

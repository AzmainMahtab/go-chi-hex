// Package ports
// ports contain all the interfaces needed for inbound
// and outbound
package ports

import (
	"context"

	"github.com/AzmainMahtab/docpad/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error

	// Read active users take optional filtering return: list of users and error if any
	Read(ctx context.Context, filter map[string]any, showDeleted bool) ([]*domain.User, error)

	// Read a single active user
	ReadOne(ctx context.Context, id int) (*domain.User, error)

	// Update provided fields and update pertially
	Update(ctx context.Context, id int, updates map[string]any) error

	// Soft delete a user set deleted_at current and user_status = 'inactive'
	SoftDelete(ctx context.Context, id int) error

	// Read soft deleted users with optional filtering
	Trash(ctx context.Context, filter map[string]any) ([]*domain.User, error)

	// Permanently delet a user
	Prune(ctx context.Context, id int) error
}

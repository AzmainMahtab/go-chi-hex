// Package postgres
// This one holds the user repository
package postgres

import (
	"context"
	"database/sql"

	"github.com/AzmainMahtab/docpad/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO "user" (user_name, email, phone, password)
		VALUSE ($1, $2, $3, $4)
		RETURNING id, user_status, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		u.UserName,
		u.Email,
		u.Phone,
		u.Password,
	).Scan(&u.ID, &u.UserStatus, &u.CreatedAt, &u.UpdatedAt)
}

// Package postgres
// This one holds the user repository
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
		/* SQL */
		INSERT INTO "user" (user_name, email, phone, password)
		VALUES ($1, $2, $3, $4)
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

func (r *UserRepo) Read(ctx context.Context, filter map[string]any, showDeleted bool) ([]*domain.User, error) {
	return nil, nil
}

func (r *UserRepo) ReadOne(ctx context.Context, id int) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT id, user_name, email, phone, user_status, created_at, updated_at
		FROM "user"
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
		&u.Phone,
		&u.UserStatus,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Record not found %v", err)
		return nil, err
	}

	return u, err
}

func (r *UserRepo) Update(ctx context.Context, id int, updates map[string]any) error {
	return nil
}

func (r *UserRepo) SoftDelete(ctx context.Context, id int) error {
	return nil
}

func (r *UserRepo) Trash(ctx context.Context, filter map[string]any) ([]*domain.User, error) {
	return nil, nil
}

func (r *UserRepo) Prune(ctx context.Context, id int) error {
	return nil
}

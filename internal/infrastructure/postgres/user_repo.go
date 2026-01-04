// Package postgres
// This one holds the user repository
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AzmainMahtab/docpad/internal/domain"
	"github.com/jmoiron/sqlx"
	"strings"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	// Wrap the standard *sql.DB into sqlx.DB
	return &UserRepo{
		db: sqlx.NewDb(db, "pgx"),
	}
}

// Create() creates a user entity
func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO "user" (user_name, email, phone, password)
		VALUES (:user_name, :email, :phone, :password)
		RETURNING id, user_status, created_at, updated_at`

	// NamedQueryContext maps :user_name to u.UserName via tags
	rows, err := r.db.NamedQueryContext(ctx, query, u)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// StructScan fills the ID and Status back into our pointer 'u'
		return rows.StructScan(u)
	}
	return rows.Err()
}

// ReadOne() reads an user entity with it's id
func (r *UserRepo) ReadOne(ctx context.Context, id int) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT * FROM "user" WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, u, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a specific domain error
		}
		return nil, err
	}
	return u, nil
}

// ReadAll() reads all the user entities with deleted users (optional)
func (r *UserRepo) ReadAll(ctx context.Context, filter map[string]any, showDeleted bool) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT * FROM "user" WHERE 1=1`

	if !showDeleted {
		query += " AND deleted_at IS NULL"
	}

	// using helper to build the query and get args
	finalQuery, args := r.appendFilters(query, filter)

	// SelectContext handles the loop and scanning for us
	err := r.db.SelectContext(ctx, &users, finalQuery, args...)
	return users, err
}

// Update() updates an user entity
func (r *UserRepo) Update(ctx context.Context, id int, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}

	// Adding the ID to the map for the WHERE clause
	updates["id"] = id

	var sb strings.Builder
	sb.WriteString(`UPDATE "user" SET `)

	// Buildng the SET string: "user_name = :user_name, email = :email"
	for col := range updates {
		if col == "id" {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s = :%s, ", col, col))
	}

	// Finalize query string
	query := strings.TrimSuffix(sb.String(), ", ")
	query += " WHERE id = :id AND deleted_at IS NULL"

	// NamedExec is great for maps! It matches map keys to :placeholders
	_, err := r.db.NamedExecContext(ctx, query, updates)
	return err
}

// SoftDelete() soft delets an user with status set to inactive and deleted_at date
func (r *UserRepo) SoftDelete(ctx context.Context, id int) error {
	query := `UPDATE "user" SET deleted_at = NOW(), user_status = 'inactive' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Trash() reads all the deletedusers
func (r *UserRepo) Trash(ctx context.Context, filter map[string]any) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT * FROM "user" WHERE deleted_at IS NOT NULL`

	finalQuery, args := r.appendFilters(query, filter)

	err := r.db.SelectContext(ctx, &users, finalQuery, args...)
	return users, err
}

// Prune() hard deletes an user
func (r *UserRepo) Prune(ctx context.Context, id int) error {
	query := `DELETE FROM "user" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// HELPER FUCTIONS
// appendFilters() helps you appednd the filters and build quries clause
func (r *UserRepo) appendFilters(baseQuery string, filter map[string]any) (string, []any) {
	var args []any
	counter := strings.Count(baseQuery, "$") + 1

	for column, value := range filter {
		baseQuery += fmt.Sprintf(" AND %s = $%d", column, counter)
		args = append(args, value)
		counter++
	}
	return baseQuery, args
}

// Package postgres
// This one holds the user repository
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (r *UserRepo) ReadAll(ctx context.Context, filter map[string]any, showDeleted bool) ([]*domain.User, error) {
	query := `
		SELECT id, user_name, email, phone, user_status, created_at, updated_at, deleted_at 
		FROM "user" 
		WHERE 1=1
	`

	// Add Soft-Delete filter
	if !showDeleted {
		query += " AND deleted_at IS NULL"
	}

	// Dynamically build the WHERE clause from the map
	query, args := r.appendFilters(query, filter)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // ALWAYS close rows to prevent memory leaks

	return r.scanRows(rows)
}

func (r *UserRepo) ReadOne(ctx context.Context, id int) (*domain.User, error) {
	u := &domain.User{}
	query := `
		/* SQL */
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
	query := `
		SELECT id, user_name, email, phone, user_status, created_at, updated_at, deleted_at
		FROM "user"
		WHERE deleted_at IS NOT NULL
	`
	query, args := r.appendFilters(query, filter)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *UserRepo) Prune(ctx context.Context, id int) error {
	return nil
}

// Itterates over the rows to send domain object
func (r *UserRepo) scanRows(rows *sql.Rows) ([]*domain.User, error) {
	var users []*domain.User

	for rows.Next() {
		var user domain.User

		// must list these in the EXACT order they appear in SELECT statement
		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.Email,
			&user.Phone,
			&user.UserStatus,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}

		users = append(users, &user)
	}

	// check for errors after the loop finishes
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return users, nil
}

// appendFilters helps to build dynamic WHERE clauses
func (r *UserRepo) appendFilters(baseQuery string, filter map[string]any) (string, []any) {
	var args []any

	//  Start counting at 1 (Postgres uses $1, $2, etc.)
	counter := 1

	//  Loop through every item in your filter map
	for column, value := range filter {
		// Build the string: " AND username = $1"
		// Then " AND email = $2", etc.
		baseQuery = baseQuery + fmt.Sprintf(" AND %s = $%d", column, counter)

		// Add the actual value (e.g., "John") to our list
		args = append(args, value)

		// Increment the number for the next loop
		counter++
	}

	return baseQuery, args
}

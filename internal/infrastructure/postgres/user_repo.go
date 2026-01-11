// Package postgres
// This one holds the user repository
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/jmoiron/sqlx"
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
		return MapError(err)
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
			return nil, MapError(err) // Or a specific domain error
		}
		return nil, MapError(err)
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

	return users, MapError(err)
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

	return MapError(err)
}

// SoftDelete() soft delets an user with status set to inactive and deleted_at date
func (r *UserRepo) SoftDelete(ctx context.Context, id int) error {
	query := `UPDATE "user" SET deleted_at = NOW(), user_status = 'inactive' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	return MapError(err)
}

// Restore() restores a trashed user
func (r *UserRepo) Restore(ctx context.Context, id int) error {
	query := `UPDATE "user" SET deleted_at = NULL, updated_at = NOW(), user_status = 'active' WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	return MapError(err)
}

// Trash() reads all the deletedusers
func (r *UserRepo) Trash(ctx context.Context, filter map[string]any) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT * FROM "user" WHERE deleted_at IS NOT NULL`

	finalQuery, args := r.appendFilters(query, filter)

	err := r.db.SelectContext(ctx, &users, finalQuery, args...)
	return users, MapError(err)
}

// Prune() hard deletes an user
func (r *UserRepo) Prune(ctx context.Context, id int) error {
	query := `DELETE FROM "user" WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return MapError(err)
}

func (r *UserRepo) CheckConflict(ctx context.Context, username, email, phone string) ([]domain.ErrorItem, error) {
	query := `
		SELECT 
			EXISTS(SELECT 1 FROM "user" WHERE user_name = $1 AND deleted_at IS NULL) as username_taken,
			EXISTS(SELECT 1 FROM "user" WHERE email = $2 AND deleted_at IS NULL) as email_taken,
			EXISTS(SELECT 1 FROM "user" WHERE phone = $3 AND deleted_at IS NULL) as phone_taken
	`

	var res struct {
		UsernameTaken bool `db:"username_taken"`
		EmailTaken    bool `db:"email_taken"`
		PhoneTaken    bool `db:"phone_taken"`
	}

	if err := r.db.GetContext(ctx, &res, query, username, email, phone); err != nil {
		return nil, MapError(err)
	}

	var conflicts []domain.ErrorItem
	if res.UsernameTaken {
		conflicts = append(conflicts, domain.ErrorItem{Field: "user_name", Message: "username already taken"})
	}
	if res.EmailTaken {
		conflicts = append(conflicts, domain.ErrorItem{Field: "email", Message: "email already registered"})
	}
	if res.PhoneTaken {
		conflicts = append(conflicts, domain.ErrorItem{Field: "phone", Message: "phone number in use"})
	}

	return conflicts, nil
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

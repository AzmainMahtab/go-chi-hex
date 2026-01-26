// Package postgres
// This one holds the user repository
package postgres

import (
	"context"
	"database/sql"

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
		INSERT INTO "user" (uuid,user_name, email,user_role, phone, password)
		VALUES (:uuid, :user_name, :email, :user_role, :phone, :password)
		RETURNING id, user_status, created_at, updated_at
	`

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
func (r *UserRepo) ReadOne(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT * FROM "user" WHERE uuid = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, u, query, id)
	if err != nil {
		return nil, MapError(err)
	}
	return u, nil
}

func (r *UserRepo) ReadByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT * FROM "user" WHERE email = $1 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, u, query, email)
	if err != nil {
		return nil, MapError(err)
	}

	return u, nil
}

// ReadAll() reads all the user entities with deleted users (optional)
func (r *UserRepo) ReadAll(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	var users []*domain.User

	//  Start with the base query
	query := `SELECT id, uuid, user_name, email, phone, user_status, created_at, updated_at 
              FROM "user" WHERE 1=1`

	//  Named arguments map for sqlx
	args := make(map[string]any)

	//  Apply Soft Delete filter
	if !filter.ShowDeleted {
		query += ` AND deleted_at IS NULL`
	}

	//  Build Dynamic Filters
	if filter.UserName != "" {
		query += ` AND user_name ILIKE :user_name`
		args["user_name"] = "%" + filter.UserName + "%"
	}

	if filter.Email != "" {
		query += ` AND email = :email`
		args["email"] = filter.Email
	}

	if filter.Phone != "" {
		query += ` AND phone = :phone`
		args["phone"] = filter.Phone
	}

	if filter.UserStatus != "" {
		query += ` AND user_status = :user_status`
		args["user_status"] = filter.UserStatus
	}

	//  Apply Pagination
	if filter.Limit > 0 {
		query += ` LIMIT :limit`
		args["limit"] = filter.Limit
	}
	if filter.Offset > 0 {
		query += ` OFFSET :offset`
		args["offset"] = filter.Offset
	}

	//  Execute using NamedQuery
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, MapError(err)
	}
	defer rows.Close()

	//  Scan results into the domain slice
	for rows.Next() {
		u := &domain.User{}
		if err := rows.StructScan(u); err != nil {
			return nil, MapError(err)
		}
		users = append(users, u)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, MapError(err)
	}

	return users, nil
}

// Update() updates an user entity
func (r *UserRepo) Update(ctx context.Context, up domain.UserUpdate) error {
	query := `
        UPDATE "user" 
        SET 
            user_name = COALESCE(:user_name, user_name),
            email = COALESCE(:email, email),
            phone = COALESCE(:phone, phone),
            user_status = COALESCE(:user_status, user_status),
            updated_at = NOW()
        WHERE id = :uuid AND deleted_at IS NULL`

	_, err := r.db.NamedExecContext(ctx, query, map[string]any{
		"id":          up.ID,
		"user_name":   up.UserName,
		"email":       up.Email,
		"phone":       up.Phone,
		"user_status": up.Status,
	})

	return MapError(err)
}

// SoftDelete() soft delets an user with status set to inactive and deleted_at date
func (r *UserRepo) SoftDelete(ctx context.Context, id string) error {
	query := `UPDATE "user" SET deleted_at = NOW(), user_status = 'inactive' WHERE uuid = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	return MapError(err)
}

// Restore() restores a trashed user
func (r *UserRepo) Restore(ctx context.Context, id string) error {
	query := `UPDATE "user" SET deleted_at = NULL, updated_at = NOW(), user_status = 'active' WHERE uuid = $1`
	_, err := r.db.ExecContext(ctx, query, id)

	return MapError(err)
}

// Trash() reads all the deletedusers
func (r *UserRepo) Trash(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	var users []*domain.User

	//  Base Query - Note the IS NOT NULL constraint to keep the trash "secret"
	query := `SELECT id, user_name, email, phone, user_status, created_at, updated_at, deleted_at 
              FROM "user" 
              WHERE deleted_at IS NOT NULL`

	args := make(map[string]any)

	//  Build Dynamic Filters (Same logic as ReadAll, but restricted to Trash)
	if filter.UserName != "" {
		query += ` AND user_name ILIKE :user_name`
		args["user_name"] = "%" + filter.UserName + "%"
	}

	if filter.Email != "" {
		query += ` AND email = :email`
		args["email"] = filter.Email
	}

	//  Apply Pagination (Always essential for Admin views with lots of data)
	if filter.Limit > 0 {
		query += ` LIMIT :limit`
		args["limit"] = filter.Limit
	}
	if filter.Offset > 0 {
		query += ` OFFSET :offstring`
		args["offset"] = filter.Offset
	}

	//  Execution
	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, MapError(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &domain.User{}
		if err := rows.StructScan(u); err != nil {
			return nil, MapError(err)
		}
		// TEMP DEBUG LOG
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepo) ReadOneDeleted(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT * FROM "user" WHERE uuid = $1 AND deleted_at IS NOT NULL`
	u := &domain.User{}

	err := r.db.GetContext(ctx, u, query, id)
	if err != nil {
		return nil, MapError(err)
	}

	return u, nil

}

// Prune() hard deletes an user
func (r *UserRepo) Prune(ctx context.Context, id string) error {
	query := `DELETE FROM "user" WHERE uuid = $1`
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

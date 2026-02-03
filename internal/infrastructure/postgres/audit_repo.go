// Package postgres
// Audit repository implementation using PostgreSQL
package postgres

import (
	"context"
	"database/sql"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/jmoiron/sqlx"
)

type AuditRepo struct {
	db *sqlx.DB
}

func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{
		db: sqlx.NewDb(db, "pgx"),
	}
}

func (r *AuditRepo) Create(ctx context.Context, auditLog domain.Audit) error {
	query := `INSERT INTO audit_log 
						(uuid, event_type, actor_id, payload) VALUES 
						(:uuid, :event_type, :actor_id, :payload)`

	_, err := r.db.NamedExecContext(ctx, query, auditLog)

	return err
}

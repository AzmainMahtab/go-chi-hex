package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// Global regex to extract field names from constraints like "users_email_key" or "appts_date_key"
var constraintRegex = regexp.MustCompile(`unique constraint ".*_(.*)_key"`)

// MapError is the global translator for any Postgres repository
func MapError(err error) error {
	if err == nil {
		return nil
	}

	//  Handle "Not Found"
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.AppError{
			Code:    domain.CodeNotFound,
			Message: "the requested resource was not found",
			Err:     err,
		}
	}

	//  Handle Postgres specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {

		slog.Error("Database execution failed",
			"code", pgErr.Code,
			"message", pgErr.Message,
			"detail", pgErr.Detail,
			"hint", pgErr.Hint,
			"where", pgErr.Where,
			"table", pgErr.TableName,
		)

		switch pgErr.Code {
		case "42703": // undefined_column
			return &domain.AppError{
				Code:    domain.CodeInternal,
				Message: fmt.Sprintf("schema mismatch: %s", pgErr.Message),
				Err:     err,
			}
		case "23505": // unique_violation
			return handleConflict(pgErr)
		case "23503": // foreign_key_violation
			return &domain.AppError{
				Code:    domain.CodeValidation,
				Message: "this operation references a record that does not exist",
				Err:     err,
			}
		case "23502": // not_null_violation
			return &domain.AppError{
				Code:    domain.CodeValidation,
				Message: "missing required information",
				Err:     err,
			}
		}
	}

	slog.Error("DATA LAYER ERR", ":", err)
	//  Default Fallback
	return &domain.AppError{
		Code:    domain.CodeInternal,
		Message: "a data layer error occurred",
		Err:     err,
	}
}

func handleConflict(pgErr *pgconn.PgError) error {
	// Extract field name dynamically
	match := constraintRegex.FindStringSubmatch(pgErr.Message)
	field := "resource"
	if len(match) > 1 {
		field = match[1]
	}

	return &domain.AppError{
		Code:    domain.CodeConflict,
		Message: fmt.Sprintf("%s is already in use", field),
		Field:   field,
		Err:     pgErr,
	}
}

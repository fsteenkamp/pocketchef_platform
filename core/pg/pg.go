package pg

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ERR_FOREIGN_KEY_VIOLATION = "23503"
	ERR_UNIQUE_VIOLATION      = "23505"
)

func IsErrForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == ERR_FOREIGN_KEY_VIOLATION
	}
	return false
}

func IsErrUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == ERR_UNIQUE_VIOLATION
	}
	return false
}

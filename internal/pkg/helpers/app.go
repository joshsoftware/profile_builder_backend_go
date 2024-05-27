package helpers

import (
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// IsDuplicateKeyError returns true if the given key is duplicate and false otherwise
func IsDuplicateKeyError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return strings.Contains(pgErr.Message, "duplicate key value violates unique constraint")
	}
	return false
}

// IsInvalidProfileError returns true if the given profile_id is wrong and false otherwise
func IsInvalidProfileError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return strings.Contains(pgErr.Message, "violates foreign key constraint")
	}
	return false
}

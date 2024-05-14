package helpers

import (
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsDuplicateKeyError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return strings.Contains(pgErr.Message, "duplicate key value violates unique constraint")
	}
	return false
}
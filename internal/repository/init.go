package repository

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// InitializeDatabase function used to initialize the database and returns the database object
func InitializeDatabase(ctx context.Context) (*pgx.Conn, error) {
	db, err := pgx.Connect(ctx, os.Getenv("PSQL_INFO"))
	if err != nil {
		return nil, errors.ErrConnectionFailed
	}

	err = db.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return db, nil
}

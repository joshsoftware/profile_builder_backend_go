package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

// InitializeDatabase function used to initialize the database and returns the database object
func InitializeDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	connPool, err := pgxpool.NewWithConfig(ctx, Config())
	if err != nil {
		zap.S().Error("Failed to create connection pool : ", err)
		return nil, errors.ErrConnectionFailed
	}

	err = connPool.Ping(ctx)
	if err != nil {
		zap.S().Error("Failed to ping the database : ", err)
		return nil, errors.ErrConnectionFailed
	}

	return connPool, nil
}

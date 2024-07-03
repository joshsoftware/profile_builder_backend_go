package repository

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"go.uber.org/zap"
)

// Confuguration for the database
func Config() *pgxpool.Config {
	const defaultMaxConns = int32(10)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(os.Getenv("PSQL_INFO"))
	if err != nil {
		zap.S().Error("Failed to configure : ", err)
		return nil
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.MaxConnLifetimeJitter = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		log.Println("Acquiring connection")
		return true
	}

	dbConfig.AfterRelease = func(conn *pgx.Conn) bool {
		log.Println("Releasing connection")
		return true
	}

	dbConfig.BeforeClose = func(conn *pgx.Conn) {
		log.Println("Closing connection")
	}
	return dbConfig
}

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

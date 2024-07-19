package repository

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

// Config defines configuration of pgxpool
func Config() *pgxpool.Config {

	dbConfig, err := pgxpool.ParseConfig(os.Getenv("PSQL_INFO"))
	if err != nil {
		zap.S().Error("Failed to configure : ", err)
		return nil
	}

	// setting the configuration
	setConfig(dbConfig)

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

func setConfig(dbConfig *pgxpool.Config) {
	dbConfig.MaxConns = helpers.ConvertStringToIntWithDefault("MAX_CONNECTIONS", constants.DefaultMaxConnections)
	dbConfig.MinConns = helpers.ConvertStringToIntWithDefault("MIN_CONNECTIONS", constants.DefaultMinConnections)
	dbConfig.MaxConnLifetime = helpers.ConvertStringToTimeDuration("MAX_CONNECTIONS_LIFETIME_IN_MINUTES", constants.DefaultConnLifeTime)
	dbConfig.MaxConnIdleTime = helpers.ConvertStringToTimeDuration("MAX_CONNECTIONS_IDLE_TIME_IN_MINUTES", constants.DefaultConnIdleTime)
	dbConfig.HealthCheckPeriod = helpers.ConvertStringToTimeDuration("HEALTH_CHECK_PERIOD_IN_MINUTES", constants.DefaultHealthCheck)
	dbConfig.MaxConnLifetimeJitter = helpers.ConvertStringToTimeDuration("CONNECT_TIMEOUT_IN_SECONDS", constants.DefaultConnectTimeout)
}

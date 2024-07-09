package repository

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"go.uber.org/zap"
)

func Config() *pgxpool.Config {
	maxConnections, err := helpers.ConvertStringToInt(os.Getenv("MAX_CONNECTIONS"))
	if err != nil {
		zap.S().Error("Failed to convert MAX_CONNECTIONS : ", err)
		return nil
	}

	minConnections, err := helpers.ConvertStringToInt(os.Getenv("MIN_CONNECTIONS"))
	if err != nil {
		zap.S().Error("Failed to convert MIN_CONNECTIONS : ", err)
		return nil
	}

	connLifeTime, err := helpers.ConvertStringToInt(os.Getenv("MAX_CONNECTIONS_LIFETIME_IN_MINUTES"))
	if err != nil {
		zap.S().Error("Failed to convert MAX_CONNECTIONS_LIFETIME_IN_MINUTES : ", err)
		return nil
	}

	connIdleTime, err := helpers.ConvertStringToInt(os.Getenv("MAX_CONNECTIONS_IDLE_TIME_IN_MINUTES"))
	if err != nil {
		zap.S().Error("Failed to convert MAX_CONNECTIONS_IDLE_TIME_IN_MINUTES : ", err)
		return nil
	}

	healthCheckPeriod, err := helpers.ConvertStringToInt(os.Getenv("HEALTH_CHECK_PERIOD_IN_MINUTES"))
	if err != nil {
		zap.S().Error("Failed to convert HEALTH_CHECK_PERIOD_IN_MINUTES : ", err)
		return nil
	}

	connectTimeout, err := helpers.ConvertStringToInt(os.Getenv("CONNECT_TIMEOUT_IN_SECONDS"))
	if err != nil {
		zap.S().Error("Failed to convert CONNECT_TIMEOUT_IN_SECONDS : ", err)
		return nil
	}

	connLifeTimeInMinute := time.Minute * time.Duration(connLifeTime)
	connIdleTimeInMinute := time.Minute * time.Duration(connIdleTime)
	healthCheckPeriodInMinute := time.Minute * time.Duration(healthCheckPeriod)
	connectTimeoutInSecond := time.Second * time.Duration(connectTimeout)

	dbConfig, err := pgxpool.ParseConfig(os.Getenv("PSQL_INFO"))
	if err != nil {
		zap.S().Error("Failed to configure : ", err)
		return nil
	}

	dbConfig.MaxConns = int32(maxConnections)
	dbConfig.MinConns = int32(minConnections)
	dbConfig.MaxConnLifetime = connLifeTimeInMinute
	dbConfig.MaxConnIdleTime = connIdleTimeInMinute
	dbConfig.HealthCheckPeriod = healthCheckPeriodInMinute
	dbConfig.MaxConnLifetimeJitter = connectTimeoutInSecond

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

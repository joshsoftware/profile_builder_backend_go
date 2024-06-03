package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

func main() {
	dbURL := os.Getenv("PSQL_INFO")
	if dbURL == "" {
		zap.S().Info("PSQL_INFO environment variable is not set")
	}

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		zap.S().Errorf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	err = runMigrations(conn)
	if err != nil {
		zap.S().Errorf("Failed to run migrations: %v\n", err)
	}

	zap.S().Info("Migrations ran successfully!")
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("DB_MIGRATE_PATH"),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	return nil
}

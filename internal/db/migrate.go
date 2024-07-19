package main

import (
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	// mainMigrationsDIR defines the directory where all migration files are located
	mainMigrationsDIR = "./internal/db/migrations"

	// mainMigrationFilesPath defines path for migration files
	mainMigrationFilesPath = "file://" + mainMigrationsDIR
)

// Migration used to define migrations
type Migration struct {
	m *migrate.Migrate
}

// InitMainDBMigrations used to initialize migrations
func InitMainDBMigrations() Migration {
	var (
		dbConnection string
		err          error
	)

	dbConnection = os.Getenv("DB_MIGRATION")
	if dbConnection == "" {
		zap.S().Fatal("PSQL_INFO environment variable is not set")
	}

	m, err := migrate.New(mainMigrationFilesPath, dbConnection)
	if err != nil {
		zap.S().Fatal("Error initializing migrations:", err)
	}

	return Migration{m}
}

// RunMigrations used to run a migrations
func (migration Migration) RunMigrations() {
	zap.S().Infof("Migrations started from %s", mainMigrationsDIR)
	startTime := time.Now()
	defer func() {
		zap.S().Infof("Migrations complete, total time taken %s", time.Since(startTime))
	}()

	dbVersion, dirty, err := migration.m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		zap.S().Fatal(err)
	}

	localVersion := uint(1)

	if dbVersion > localVersion {
		zap.S().Fatalf("Your database migration %d is ahead of local migration %d, you might need to rollback a few migrations", dbVersion, localVersion)
	}
	if dbVersion < localVersion && dirty {
		zap.S().Fatalf("Your currently active database migration %d is dirty, you might need to rollback it and then deploy again.", dbVersion)
	}

	err = migration.m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return
		}

		dbVersion, _, err2 := migration.m.Version()
		if err2 != nil {
			zap.S().Fatal(err2)
		}

		zap.S().Fatalf("Migration failed with error %s, current active database migration is %d", err, dbVersion)
	}
}

// MigrationsUp used to make migrations up
func (migration Migration) MigrationsUp() {
	err := migration.m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			zap.S().Error("No new migrations to apply")
			return
		}
		zap.S().Fatal(err)
		return
	}
	zap.S().Info("Migration up complete")
}

// MigrationsDown used to make migrations down
func (migration Migration) MigrationsDown() {
	err := migration.m.Down()
	if err != nil {
		if err == migrate.ErrNoChange {
			zap.S().Info("No migrations to revert")
			return
		}

		zap.S().Fatal(err)
		return
	}

	zap.S().Info("Migration down complete")
}

func main() {
	//setting logger
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// Setup env
	err := godotenv.Load()
	if err != nil {
		zap.S().Info("Error loading .env file:", err)
		return
	}

	if len(os.Args) < 2 {
		zap.S().Error("Missing action argument. Use 'up' or 'down'.")
		os.Exit(1)
	}

	migration := InitMainDBMigrations()

	action := os.Args[1]
	switch action {
	case "up":
		migration.MigrationsUp()
	case "down":
		migration.MigrationsDown()
	default:
		zap.S().Info("Invalid action. Use 'up' or 'down'.")
	}
}

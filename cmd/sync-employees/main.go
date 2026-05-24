package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	intranetclient "github.com/joshsoftware/profile_builder_backend_go/internal/client/intranet"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/log"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Set up zap logger
	logger, err := log.SetupLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Load .env
	if err := godotenv.Load(); err != nil {
		zap.S().Warn("No .env file found, relying on environment variables")
	}

	// Validate required env vars
	baseURL := os.Getenv("INTRANET_API_BASE_URL")
	apiKey := os.Getenv("INTRANET_API_KEY")
	if baseURL == "" || apiKey == "" {
		zap.S().Error("INTRANET_API_BASE_URL and INTRANET_API_KEY must be set")
		os.Exit(1)
	}

	// Initialize DB
	db, err := repository.InitializeDatabase(ctx)
	if err != nil {
		zap.S().Error("Failed to connect to database: ", err)
		os.Exit(1)
	}
	defer db.Close()
	zap.S().Info("Connected to database")

	// Build dependencies
	repoDeps := service.RepoDeps{
		ProfileDeps:    repository.NewProfileRepo(db),
		IntranetClient: intranetclient.NewClient(baseURL, apiKey),
	}

	svc := service.NewServices(repoDeps)

	// Run sync
	zap.S().Info("Starting employee ID sync...")
	updated, skipped, err := svc.SyncEmployees(ctx)
	if err != nil {
		zap.S().Error("Employee sync failed: ", err)
		os.Exit(1)
	}

	summary := fmt.Sprintf("Sync complete. Updated: %d, Skipped: %d", updated, skipped)
	zap.S().Info(summary)
	fmt.Println(summary)
}

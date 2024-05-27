package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	//setting logger
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	//setup env
	err := godotenv.Load()
	if err != nil {
		logger.Info("error loading.env file")
		return
	}

	fmt.Println("Starting Server...")
	defer fmt.Println("Shutting Down Server...")

	//Initialize DB
	db, err := repository.InitializeDatabase(ctx)
	if err != nil {
		logger.Error("Database Error  : ", zap.Error(err))
		return
	}
	fmt.Println("Connected to Database!")

	//Initializing Services
	services := service.NewServices(db)

	//Initializing Router
	router := api.NewRouter(ctx, services)

	// CORS middleware
	cors := cors.New(constants.CorsOptions)

	// Setup the server
	server := &http.Server{
		Addr:    os.Getenv("PORT_INFO"),
		Handler: cors.Handler(router),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error starting server: ", zap.Error(err))
		}
	}()

	//Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)

	sig := <-signalChan
	fmt.Println("Received terminate, gracefully shutting down: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
	cancel()

	//Shutdown any other resources here, e.g., database connections
	if err := db.Close(ctx); err != nil {
		logger.Error("Database close error: ", zap.Error(err))
	}
}

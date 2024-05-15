package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
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

	//Creating Services
	services := app.NewServices(ctx, db)

	//Initializaing Router
	router := api.NewRouter(ctx, services)

	err = http.ListenAndServe("localhost:1925", router)
	if err != nil {
		logger.Error("Error Starting Server  : ", zap.Error(err))
		return
	}
}

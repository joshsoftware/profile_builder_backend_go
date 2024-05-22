package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	//setup env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading.env file")
		return
	}

	//setting logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	fmt.Println("Starting Server...")
	defer fmt.Println("Shutting Down Server...")

	//Initialize DB
	db, err := repository.InitializeDatabase(ctx)
	if err != nil {
		fmt.Println("Database Error : ", err)
		logger.Error("Database Error  : ", zap.Error(err))
		return
	}
	fmt.Println("Connected to Database!")

	//Creating Services
	services := app.NewServices(db, ctx)

	//Initializaing Router
	router := api.NewRouter(services, ctx)
	cors := cors.New(constants.CorsOptions)

	err = http.ListenAndServe("localhost:1925", cors.Handler(router))
	if err != nil {
		fmt.Println("Error Starting Server : ", err)
		logger.Error("Error Starting Server  : ", zap.Error(err))
		return
	}
}

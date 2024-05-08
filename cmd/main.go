package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

func main() {
	ctx := context.Background();

	//setup env
	err := godotenv.Load();
	if err != nil {
		fmt.Println("error loading.env file");
        return
	}

	fmt.Println("Starting Server...");
	defer fmt.Println("Shutting Down Server...");

	//Initialize DB
	db, err := repository.InitializeDatabase(ctx);
	if err != nil{
		fmt.Println("Database Error : ", err);
		return 
	}
	fmt.Println("Connected to Database!")

	//Creating Services
	services := app.NewServices(db);

	//Initializaing Router
	router := api.NewRouter(services);

	err = http.ListenAndServe("localhost:1925", router)
	if err != nil{
        fmt.Println("Error Starting Server : ", err);
        return
    }


}
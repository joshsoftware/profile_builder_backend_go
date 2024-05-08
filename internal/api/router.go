package api

import (
	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
)

func NewRouter(deps app.Dependencies) *mux.Router{
	router := mux.NewRouter();
	
	return router
}
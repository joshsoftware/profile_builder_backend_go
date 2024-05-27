package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
)

// NewRouter returns a object that contains all routes of application
func NewRouter(ctx context.Context, svc service.Service) *mux.Router {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/profiles").Subrouter()

	// api/profiles

	//POST APIs
	router.HandleFunc("/profiles", handler.CreateProfileHandler(ctx, svc)).Methods(http.MethodPost)
	subrouter.HandleFunc("/educations", handler.CreateEducationHandler(ctx, svc)).Methods(http.MethodPost)

	//GET APIs
	return router
}

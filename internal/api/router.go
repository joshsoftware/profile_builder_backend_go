package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	get "github.com/joshsoftware/profile_builder_backend_go/internal/api/GET"
	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
)

func NewRouter(ctx context.Context, deps app.Dependencies) *mux.Router {
	router := mux.NewRouter()

	//POST APIs
	router.HandleFunc("/profiles", post.CreateProfileHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)

	subrouter := router.PathPrefix("/profiles").Subrouter()
	subrouter.HandleFunc("/educations", post.CreateEducationHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)
	subrouter.HandleFunc("/projects", post.CreateProjectsHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)

	//GET APIs
	router.HandleFunc("/list_profiles", get.GetProfileListHandler(ctx, deps.ProfileService)).Methods(http.MethodGet)

	return router
}

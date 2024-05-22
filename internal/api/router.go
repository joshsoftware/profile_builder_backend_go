package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	get "github.com/joshsoftware/profile_builder_backend_go/internal/api/GET"
	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
)

// NewRouter returns a object that contains all routes of application
func NewRouter(ctx context.Context, deps app.Dependencies) *mux.Router {
	router := mux.NewRouter()

	//POST APIs
	router.HandleFunc("/profiles", post.CreateProfileHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)

	subrouter := router.PathPrefix("/profiles").Subrouter()
	subrouter.HandleFunc("/educations", post.CreateEducationHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)
	subrouter.HandleFunc("/projects", post.CreateProjectHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)
	subrouter.HandleFunc("/experiences", post.CreateExperienceHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)
	subrouter.HandleFunc("/certificates", post.CreateCertificateHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)
	subrouter.HandleFunc("/achievements", post.CreateAchievementHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)

	//GET APIs
	router.HandleFunc("/list_profiles", get.ProfileListHandler(ctx, deps.ProfileService)).Methods(http.MethodGet)

	return router
}

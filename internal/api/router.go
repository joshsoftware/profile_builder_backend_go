package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
)

func NewRouter(ctx context.Context, deps app.Dependencies) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/profiles", post.CreateProfileHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)

	//Subrouter routes of profiles
	subrouter := router.PathPrefix("/profiles").Subrouter()
	subrouter.HandleFunc("/educations", post.CreateEducationHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)
	subrouter.HandleFunc("/projects", post.CreateProjectsHandler(ctx, deps.ProfileService)).Methods(http.MethodPost)

	return router
}

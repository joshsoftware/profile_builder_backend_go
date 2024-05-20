package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app"
)

func NewRouter(deps app.Dependencies, ctx context.Context) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/profiles", post.CreateProfileHandler(deps.ProfileService, ctx)).Methods(http.MethodPost)
	router.HandleFunc("/login", post.Login(ctx, deps.ProfileService)).Methods(http.MethodPost)
	return router
}

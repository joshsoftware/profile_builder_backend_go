package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
)

// NewRouter returns a object that contains all routes of application
func NewRouter(ctx context.Context, svc service.Service) *mux.Router {
	router := mux.NewRouter()

	// user login router
	router.HandleFunc("/login", handler.Login(ctx, svc)).Methods(http.MethodPost)

	profileSubrouter := router.PathPrefix("/api").Subrouter()
	profileSubrouter.Use(middleware.AuthMiddleware)

	//POST APIs
	profileSubrouter.HandleFunc("/profiles", handler.CreateProfileHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations", handler.CreateEducationHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{project_id}/projects", handler.CreateProjectHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences", handler.CreateExperienceHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates", handler.CreateCertificateHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profiles_id}/achievements", handler.CreateAchievementHandler(ctx, svc)).Methods(http.MethodPost)

	//GET APIs
	profileSubrouter.HandleFunc("/list_profiles", handler.ProfileListHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}", handler.GetProfileHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations", handler.GetEducationHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/projects", handler.GetProjectHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences", handler.GetExperienceHandler(ctx, svc)).Methods(http.MethodGet)

	// PUT APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}", handler.UpdateProfileHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations/{id}", handler.UpdateEducationHandler(ctx, svc)).Methods(http.MethodPut)

	return router
}

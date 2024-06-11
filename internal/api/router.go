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

	// Profile APIs
	profileSubrouter.HandleFunc("/profiles", handler.CreateProfileHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}", handler.UpdateProfileHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles", handler.ProfileListHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles", handler.GetProfileHandler(ctx, svc)).Methods(http.MethodGet)

	// Educations APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations", handler.CreateEducationHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations", handler.GetEducationHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations/{id}", handler.UpdateEducationHandler(ctx, svc)).Methods(http.MethodPut)

	// Certificates APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates", handler.CreateCertificateHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates", handler.ListCertificatesHandler(ctx, svc)).Methods(http.MethodGet)

	// Projects APIs
	profileSubrouter.HandleFunc("/profiles/{project_id}/projects", handler.CreateProjectHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/projects", handler.GetProjectHandler(ctx, svc)).Methods(http.MethodGet)

	// Experiences APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences", handler.CreateExperienceHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences", handler.GetExperienceHandler(ctx, svc)).Methods(http.MethodGet)

	// Achievements APIs
	profileSubrouter.HandleFunc("/profiles/{profiles_id}/achievements", handler.CreateAchievementHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/achievements", handler.ListAchievementsHandler(ctx, svc)).Methods(http.MethodGet)

	return router
}

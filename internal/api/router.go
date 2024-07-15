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
	profileSubrouter.HandleFunc("/profiles/{profile_id}", handler.GetProfileHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/skills", handler.SkillsListHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}", handler.DeleteProfileHandler(ctx, svc)).Methods(http.MethodDelete)

	// Educations APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations", handler.CreateEducationHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations", handler.ListEducationHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations/{id}", handler.UpdateEducationHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/educations/{id}", handler.DeleteEducationHandler(ctx, svc)).Methods(http.MethodDelete)

	// Certificates APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates", handler.CreateCertificateHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates", handler.ListCertificatesHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates/{id}", handler.UpdateCertificateHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/certificates/{id}", handler.DeleteCertificatesHandler(ctx, svc)).Methods(http.MethodDelete)

	// Projects APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/projects", handler.CreateProjectHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/projects", handler.ListProjectHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/projects/{id}", handler.UpdateProjectHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/projects/{id}", handler.DeleteProjectHandler(ctx, svc)).Methods(http.MethodDelete)

	// Experiences APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences", handler.CreateExperienceHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences", handler.ListExperienceHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences/{id}", handler.UpdateExperienceHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/experiences/{id}", handler.DeleteExperienceHandler(ctx, svc)).Methods(http.MethodDelete)

	// Achievements APIs
	profileSubrouter.HandleFunc("/profiles/{profile_id}/achievements", handler.CreateAchievementHandler(ctx, svc)).Methods(http.MethodPost)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/achievements", handler.ListAchievementsHandler(ctx, svc)).Methods(http.MethodGet)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/achievements/{id}", handler.UpdateAchievementHandler(ctx, svc)).Methods(http.MethodPut)
	profileSubrouter.HandleFunc("/profiles/{profile_id}/achievements/{id}", handler.DeleteAchievementHandler(ctx, svc)).Methods(http.MethodDelete)

	return router
}

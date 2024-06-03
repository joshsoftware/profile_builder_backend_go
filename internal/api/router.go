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

	//POST APIs
	router.HandleFunc("/profiles", handler.CreateProfileHandler(ctx, svc)).Methods(http.MethodPost)
	subrouter.HandleFunc("/{profile_id}/educations", handler.CreateEducationHandler(ctx, svc)).Methods(http.MethodPost)
	subrouter.HandleFunc("/{profile_id}/projects", handler.CreateProjectHandler(ctx, svc)).Methods(http.MethodPost)
	subrouter.HandleFunc("/{profile_id}/experiences", handler.CreateExperienceHandler(ctx, svc)).Methods(http.MethodPost)
	subrouter.HandleFunc("/{profile_id}/certificates", handler.CreateCertificateHandler(ctx, svc)).Methods(http.MethodPost)
	subrouter.HandleFunc("/{profile_id}/achievements", handler.CreateAchievementHandler(ctx, svc)).Methods(http.MethodPost)

	//GET APIs
	router.HandleFunc("/profiles", handler.ProfileListHandler(ctx, svc)).Methods(http.MethodGet)
	router.HandleFunc("/skills", handler.SkillsListHandler(ctx, svc)).Methods(http.MethodGet)
	router.HandleFunc("/profiles/{profile_id}", handler.GetProfileHandler(ctx, svc)).Methods(http.MethodGet)
	subrouter.HandleFunc("/{profile_id}/educations", handler.GetEducationHandler(ctx, svc)).Methods(http.MethodGet)
	subrouter.HandleFunc("/{profile_id}/projects", handler.GetProjectHandler(ctx, svc)).Methods(http.MethodGet)
	subrouter.HandleFunc("/{profile_id}/experiences", handler.GetExperienceHandler(ctx, svc)).Methods(http.MethodGet)

	//PUT APIs
	router.HandleFunc("/profiles/{profile_id}", handler.UpdateProfileHandler(ctx, svc)).Methods(http.MethodPut)
	subrouter.HandleFunc("/{profile_id}/educations/{id}", handler.UpdateEducationHandler(ctx, svc)).Methods(http.MethodPut)
	subrouter.HandleFunc("/{profile_id}/experiences/{id}", handler.UpdateExperienceHandler(ctx, svc)).Methods(http.MethodPut)
	subrouter.HandleFunc("/{profile_id}/certificates/{id}", handler.UpdateCertificateHandler(ctx, svc)).Methods(http.MethodPut)
	subrouter.HandleFunc("/{profile_id}/projects/{id}", handler.UpdateProjectHandler(ctx, svc)).Methods(http.MethodPut)
	subrouter.HandleFunc("/{profile_id}/achievements/{id}", handler.UpdateAchievementHandler(ctx, svc)).Methods(http.MethodPut)

	return router
}

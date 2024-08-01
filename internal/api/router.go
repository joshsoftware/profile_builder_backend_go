package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
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
	profileSubrouter.Handle("/profiles", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.CreateProfileHandler(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateProfileHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.ProfileListHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.GetProfileHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/skills", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.SkillsListHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.DeleteProfileHandler(ctx, svc)))).Methods(http.MethodDelete)
	profileSubrouter.Handle("/updateSequence", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateSequenceHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles/{profile_id}", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.UpdateProfileStatusHandler(ctx, svc)))).Methods(http.MethodPatch)

	// Educations APIs
	profileSubrouter.Handle("/profiles/{profile_id}/educations", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.CreateEducationHandler(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}/educations", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.ListEducationHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}/educations/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateEducationHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles/{profile_id}/educations/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.DeleteEducationHandler(ctx, svc)))).Methods(http.MethodDelete)

	// Certificates APIs
	profileSubrouter.Handle("/profiles/{profile_id}/certificates", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.CreateCertificateHandler(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}/certificates", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.ListCertificatesHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}/certificates/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateCertificateHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles/{profile_id}/certificates/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.DeleteCertificatesHandler(ctx, svc)))).Methods(http.MethodDelete)

	// Projects APIs
	profileSubrouter.Handle("/profiles/{profile_id}/projects", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.CreateProjectHandler(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}/projects", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.ListProjectHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}/projects/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateProjectHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles/{profile_id}/projects/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.DeleteProjectHandler(ctx, svc)))).Methods(http.MethodDelete)

	// Experiences APIs
	profileSubrouter.Handle("/profiles/{profile_id}/experiences", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.CreateExperienceHandler(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}/experiences", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.ListExperienceHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}/experiences/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateExperienceHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles/{profile_id}/experiences/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.DeleteExperienceHandler(ctx, svc)))).Methods(http.MethodDelete)

	// Achievements APIs
	profileSubrouter.Handle("/profiles/{profile_id}/achievements", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.CreateAchievementHandler(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}/achievements", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.ListAchievementsHandler(ctx, svc)))).Methods(http.MethodGet)
	profileSubrouter.Handle("/profiles/{profile_id}/achievements/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.UpdateAchievementHandler(ctx, svc)))).Methods(http.MethodPut)
	profileSubrouter.Handle("/profiles/{profile_id}/achievements/{id}", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.DeleteAchievementHandler(ctx, svc)))).Methods(http.MethodDelete)

	// User Email APIs
	profileSubrouter.Handle("/profiles/{profile_id}/employee_invite", middleware.RoleMiddleware([]string{constants.Admin})(http.HandlerFunc(handler.SendUserInvitation(ctx, svc)))).Methods(http.MethodPost)
	profileSubrouter.Handle("/profiles/{profile_id}/profile_complete", middleware.RoleMiddleware([]string{constants.Admin, constants.Employee})(http.HandlerFunc(handler.SendAdminInvitation(ctx, svc)))).Methods(http.MethodPatch)

	return router
}

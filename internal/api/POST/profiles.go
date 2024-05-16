package post

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/profile"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
)

func CreateProfileHandler(ctx context.Context, profileSvc profile.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeCreateProfileRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = profileSvc.CreateProfile(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponse{
			Message: "Basic info added successfully",
		})
	}
}

func CreateEducationHandler(ctx context.Context, profileSvc profile.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeCreateEducationRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = profileSvc.CreateEducation(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponse{
			Message: "Education(s) added successfully",
		})
	}
}

func CreateProjectsHandler(ctx context.Context, profileSvc profile.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeCreateProjectRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = profileSvc.CreateProject(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponse{
			Message: "Profile(s) added successfully",
		})
	}
}

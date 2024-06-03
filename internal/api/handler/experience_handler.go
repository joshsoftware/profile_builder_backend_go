package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"go.uber.org/zap"
)

// CreateExperienceHandler handles HTTP requests to add experience details to a user profile.
func CreateExperienceHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helpers.GetParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeCreateExperinceRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		profileID, err := profileSvc.CreateExperience(ctx, req, id)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create experiences : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponseWithID{
			Message:   "Experience(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// GetExperienceHandler returns an HTTP handler that particular experiences using profileSvc.
func GetExperienceHandler(ctx context.Context, expSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		values, err := expSvc.GetExperience(ctx, profileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to get experience : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, dto.ResponseExperience{
			Experiences: values,
		})
	}
}

// UpdateExperienceHandler returns an HTTP handler that updates experience using profileSvc.
func UpdateExperienceHandler(ctx context.Context, eduSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, eduID, err := helpers.GetMultipleParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeUpdateExperienceRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		value, err := eduSvc.UpdateExperience(ctx, profileID, eduID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update experience : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, dto.MessageResponseWithID{
			Message:   "Experience updated successfully",
			ProfileID: value,
		})
	}
}

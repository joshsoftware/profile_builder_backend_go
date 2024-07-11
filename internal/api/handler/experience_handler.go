package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// CreateExperienceHandler handles HTTP requests to add experience details to a user profile.
func CreateExperienceHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
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

		profileID, err = profileSvc.CreateExperience(ctx, req, profileID, userID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create experiences : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, specs.MessageResponseWithID{
			Message:   "Experience(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// ListExperienceHandler returns an HTTP handler that particular experiences using profileSvc.
func ListExperienceHandler(ctx context.Context, expSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		filter, err := helpers.DecodeExperiencesRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrDecodeRequest)
			zap.S().Error(err)
			return
		}

		expResp, err := expSvc.ListExperiences(ctx, profileID, filter)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailespecsFetch)
			zap.S().Error("Unable to fetch experiences : ", err, "for profile id : ", profileID)
			return
		}

		if len(expResp) == 0 {
			expResp = []specs.ExperienceResponse{}
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ResponseExperience{
			Experiences: expResp,
		})
	}
}

// UpdateExperienceHandler returns an HTTP handler that updates experience using profileSvc.
func UpdateExperienceHandler(ctx context.Context, eduSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, expID, err := helpers.GetMultipleParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
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

		updatedResp, err := eduSvc.UpdateExperience(ctx, profileID, expID, userID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update experience : ", err, "for profile id : ", profileID, "experience id : ", expID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Experience updated successfully",
			ProfileID: updatedResp,
		})
	}
}

func DeleteExperienceHandler(ctx context.Context, expSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, expID, err := helpers.GetMultipleParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("error while getting the IDs from request : ", err)
			return
		}

		req := decodeDeleteExperienceRequest(profileID, expID)

		err = expSvc.DeleteExperience(ctx, req)

		if err != nil {
			if err == errors.ErrNoData {
				middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
					Message: "No data found for deletion",
				})
				return
			}
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToDelete)
			zap.S().Error("error while deleting the experience: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Experience deleted successfully",
		})
	}
}

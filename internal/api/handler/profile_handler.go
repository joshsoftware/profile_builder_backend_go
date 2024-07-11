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

// CreateProfileHandler handles HTTP requests to create a new user profile.
func CreateProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeCreateProfileRequest(r)
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

		profileID, err := profileSvc.CreateProfile(ctx, req, userID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create profile : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, specs.MessageResponseWithID{
			Message:   "Basic info added successfully",
			ProfileID: profileID,
		})
	}
}

// ProfileListHandler returns an HTTP handler that lists profiles using profileSvc.
func ProfileListHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profResponse, err := profileSvc.ListProfiles(ctx)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to list profiles : ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ListProfilesResponse{
			Profiles: profResponse,
		})
	}
}

// SkillsListHandler returns an HTTP handler that lists skills using profileSvc.
func SkillsListHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		skillsResp, err := profileSvc.ListSkills(ctx)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to list skills : ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, skillsResp)
	}
}

// GetProfileHandler returns an HTTP handler that fetches particular profile using profileSvc.
func GetProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		profResp, err := profileSvc.GetProfile(ctx, profileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to get profile : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ProfileResponse{
			Profile: profResp,
		})
	}
}

// UpdateProfileHandler returns an HTTP handler that updates profile using profileSvc.
func UpdateProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
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

		req, err := decodeUpdateProfileRequest(r)
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

		updatedResp, err := profileSvc.UpdateProfile(ctx, profileID, userID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update profile : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Basic info updated successfully",
			ProfileID: updatedResp,
		})
	}
}

// DeleteProfileHandler returns an HTTP handler that deletes profile using profileSvc.
func DeleteProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("error while getting the IDs from request : ", err)
			return
		}

		err = profileSvc.DeleteProfile(ctx, profileID)
		if err != nil {
			if err == errors.ErrNoData {
				middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
					Message: "No data found for deletion",
				})
				return
			}
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToDelete)
			zap.S().Error("error while deleting the profile: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Profile deleted successfully",
		})
	}
}

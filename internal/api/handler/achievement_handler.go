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

// CreateAchievementHandler handles HTTP requests to add achievements details to a user profile.
func CreateAchievementHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
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
		req, err := decodeCreateAchievementRequest(r)
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

		profileID, err = profileSvc.CreateAchievement(ctx, req, profileID, userID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create achievement : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, specs.MessageResponseWithID{
			Message:   "Achievement(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// UpdateAchievementHandler returns an HTTP handler that updates achievement using profileSvc.
func UpdateAchievementHandler(ctx context.Context, achSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, achID, err := helpers.GetMultipleParams(r)
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

		req, err := decodeUpdateAchievementRequest(r)
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

		updatedResp, err := achSvc.UpdateAchievement(ctx, profileID, achID, userID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update achievement : ", err, " for profile id : ", profileID, "achievement id : ", achID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Achievement updated successfully",
			ProfileID: updatedResp,
		})
	}
}

// ListAchievementsHandler returns an HTTP handler that particular achievements using profileSvc.
func ListAchievementsHandler(ctx context.Context, achSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrInvalidProfile)
			zap.S().Error(err)
			return
		}

		filter, err := helpers.DecodeAchievementRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrDecodeRequest)
			zap.S().Error(err)
			return
		}

		achivementsResp, err := achSvc.ListAchievements(ctx, profileID, filter)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailespecsFetch)
			zap.S().Error("Unable to fetch achievement : ", err, "for profile id : ", profileID)
			return
		}

		if len(achivementsResp) == 0 {
			achivementsResp = []specs.AchievementResponse{}
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ResponseAchievement{
			Achievements: achivementsResp,
		})
	}
}

// DeleteAchievementHandler returns an HTTP handler that deletes particular achievement using profileSvc.
func DeleteAchievementHandler(ctx context.Context, achSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, achievementID, err := helpers.GetMultipleParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("error while getting the IDs from request")
			return
		}

		// call the service
		err = achSvc.DeleteAchievement(ctx, profileID, achievementID)
		if err != nil {
			if err == errors.ErrNoData {
				middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
					Message: constants.NoResourceFound,
				})
				return
			}
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToDelete)
			zap.S().Error("error while deleting the achievements: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Achievement deleted successfully",
		})
	}

}

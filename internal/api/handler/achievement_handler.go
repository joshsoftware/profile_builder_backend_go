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

// CreateAchievementHandler handles HTTP requests to add achievements details to a user profile.
func CreateAchievementHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		profileID, err := profileSvc.CreateAchievement(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create achievement : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponseWithID{
			Message:   "Achievement(s) added successfully",
			ProfileID: profileID,
		})
	}
}

func ListAchievementsHandler(ctx context.Context, achSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		ID, err := helpers.ConvertStringToInt(profileID)
		if err != nil {
			zap.S().Error("error to get achievement : ", err, " for profile id : ", profileID)
			return
		}

		achivementsResp, err := achSvc.ListAchievements(ctx, ID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to fetch achievement : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, dto.ResponseAchievement{
			Achievements: achivementsResp,
		})
	}
}

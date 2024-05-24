package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
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
			zap.S().Error(err)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponseWithID{
			Message:   "Achievement(s) added successfully",
			ProfileID: profileID,
		})
	}
}

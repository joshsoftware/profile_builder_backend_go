package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// SendUserInvitation sends an invitation to the user
func SendUserInvitation(ctx context.Context, userService service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeSendUserInvitationRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("Error decoding request: ", err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidBody)
			zap.S().Error("Error validating request: ", err)
			return
		}

		profileID, err := helpers.GetProfileId(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		if profileID != req.ProfileID {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidProfile)
			zap.S().Error(errors.ErrInvalidProfile)
			return
		}

		err = userService.SendUserInvitation(ctx, userID, req)
		if err != nil {
			zap.S().Errorf("Error sending invitation: ", err)
			middleware.ErrorResponse(w, http.StatusInternalServerError, errors.ErrUnableToSendEmail)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Invitation sent successfully to employee",
		})
	}
}

func SendAdminInvitation(ctx context.Context, userService service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeSendUserInvitationRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("Error decoding request: ", err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidBody)
			zap.S().Error("Error validating request: ", err)
			return
		}

		profileID, err := helpers.GetProfileId(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		if profileID != req.ProfileID {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidProfile)
			zap.S().Error(errors.ErrInvalidProfile)
			return
		}

		err = userService.SendAdminInvitation(ctx, userID, req)
		if err != nil {
			zap.S().Errorf("Error sending invitation: %v", err)
			middleware.ErrorResponse(w, http.StatusInternalServerError, errors.ErrUnableToSendEmail)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Profile Completed Successfully",
		})
	}
}

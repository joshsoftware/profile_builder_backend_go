package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

func SendUserInvitation(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeSendUserInvitationRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("Error decoding request : ", err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("Error validating request : ", err)
			return
		}

		err = profileSvc.SendUserInvitation(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			zap.S().Error("Error sending invitation : ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Invitation sent successfully",
		})
	}
}

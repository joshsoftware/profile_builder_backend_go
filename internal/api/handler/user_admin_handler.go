package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// InviteAdmin sends an invitation to a new admin
func InviteAdmin(ctx context.Context, userService service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			zap.S().Error("Error getting user id from context: ", err)
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		var req specs.AdminInviteRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			zap.S().Error("Error decoding request body: ", err)
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidBody)
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(req.Email)

		if req.Name == "" || req.Email == "" {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidAdminRequest)
			return
		}

		err = userService.InviteAdmin(ctx, userID, req)
		if err != nil {
			if err == errors.ErrDuplicateKey {
				middleware.ErrorResponse(w, http.StatusConflict, err)
				return
			}
			zap.S().Error("Error in sending admin invitation: ", err)
			middleware.ErrorResponse(w, http.StatusInternalServerError, errors.ErrUnableToSendEmail)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Admin invited successfully",
		})
	}
}

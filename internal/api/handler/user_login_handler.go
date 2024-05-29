package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"go.uber.org/zap"
)

func Login(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userInfo dto.UserInfo
		req, err := decodeUserLoginRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("Unable to decode request : ", err)
			return
		}

		body, err := helpers.SendRequest(ctx, os.Getenv("GOOGLE_USER_INFO_URL"), req.AccessToken)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("Unable to send request to google : ", err)
			return
		}
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&userInfo); err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			zap.S().Error("Unable to decode response body : ", err)
			return
		}

		if len(userInfo.Email) == 0 {
			middleware.ErrorResponse(w, http.StatusNotFound, errors.ErrEmailNotFound)
			zap.S().Error("Invalid email")
			return
		}

		token, err := profileSvc.GenerateLoginToken(ctx, userInfo.Email)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			zap.S().Error("Unable to generate token : ", err)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)

		loginResp := dto.UserLoginResponse{
			Message:    "Login successful",
			Token:      token,
			StatusCode: http.StatusOK,
		}

		middleware.SuccessResponse(w, http.StatusOK, loginResp)
	}
}

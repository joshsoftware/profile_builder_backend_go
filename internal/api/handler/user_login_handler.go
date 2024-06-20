package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"go.uber.org/zap"
)

func Login(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userInfo specs.UserInfo
		req, err := decodeUserLoginRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(errors.ErrDecodeRequest, " : ", err)
			return
		}

		if req.AccessToken == "" {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrEmptyAccessToken)
			zap.S().Error(errors.ErrEmptyAccessToken)
			return
		}

		headers := map[string]string{
			"Content-Type": "application/json",
		}

		body, err := helpers.SendRequest(ctx, http.MethodGet, os.Getenv("GOOGLE_USER_INFO_URL"), req.AccessToken, nil, headers)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(errors.ErrGoogleRequest, ": ", err)
			return
		}

		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&userInfo); err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			zap.S().Error(errors.ErrDecodeResponse, ": ", err)
			return
		}

		if len(userInfo.Email) == 0 {
			middleware.ErrorResponse(w, http.StatusNotFound, errors.ErrEmailNotFound)
			zap.S().Error(errors.ErrInvalidEmail)
			return
		}

		token, err := profileSvc.GenerateLoginToken(ctx, userInfo.Email)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			zap.S().Error(errors.ErrGenerateToken, " : ", err)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)

		loginResp := specs.UserLoginResponse{
			Message:    "Login successful",
			Token:      token,
			StatusCode: http.StatusOK,
		}

		middleware.SuccessResponse(w, http.StatusOK, loginResp)
	}
}

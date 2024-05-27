package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
)

func Login(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeUserLoginRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		serverRequest, err := http.NewRequest("GET", os.Getenv("GOOGLE_USER_INFO_URL"), nil)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		serverRequest.Header.Set("Authorization", "Bearer "+req.AccessToken)

		resp, err := http.DefaultClient.Do(serverRequest)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		var userInfo dto.UserInfo
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&userInfo); err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		token, err := profileSvc.GenerateLoginToken(ctx, userInfo.Email)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusInternalServerError, err)
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

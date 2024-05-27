package post

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/profile"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
)

func CreateProfileHandler(profileSvc profile.Service, ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeCreateProfileRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = profileSvc.CreateProfile(req, ctx)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponse{
			Message: "Basic info added successfully",
		})
	}
}

func Login(ctx context.Context, profileSvc profile.Service) func(http.ResponseWriter, *http.Request) {
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

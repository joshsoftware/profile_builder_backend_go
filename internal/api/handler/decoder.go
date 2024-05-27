package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
)

func decodeUserLoginRequest(r *http.Request) (dto.UserLoginRequest, error) {
	var req dto.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.UserLoginRequest{}, errors.New("invalid Json in request body")
	}

	return req, nil
}

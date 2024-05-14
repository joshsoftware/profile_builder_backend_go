package post

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
)

func decodeCreateProfileRequest(r *http.Request) (dto.CreateProfileRequest, error) {
	var req dto.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateProfileRequest{}, errors.New("invalid Json in request body")
	}

	return req, nil
}
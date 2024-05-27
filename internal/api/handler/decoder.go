package handler

import (
	"encoding/json"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// Decodes the Profile Creation object Request
func decodeCreateProfileRequest(r *http.Request) (dto.CreateProfileRequest, error) {
	var req dto.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateProfileRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Updation object Request
func decodeUpdateProfileRequest(r *http.Request) (dto.UpdateProfileRequest, error) {
	var req dto.UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.UpdateProfileRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Education object Request
func decodeCreateEducationRequest(r *http.Request) (dto.CreateEducationRequest, error) {
	var req dto.CreateEducationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateEducationRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Education Updation object Request
func decodeUpdateEducationRequest(r *http.Request) (dto.UpdateEducationRequest, error) {
	var req dto.UpdateEducationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.UpdateEducationRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Project object request
func decodeCreateProjectRequest(r *http.Request) (dto.CreateProjectRequest, error) {
	var req dto.CreateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateProjectRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Experience object Request
func decodeCreateExperinceRequest(r *http.Request) (dto.CreateExperienceRequest, error) {
	var req dto.CreateExperienceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateExperienceRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Certicates object Request
func decodeCreateCertificateRequest(r *http.Request) (dto.CreateCertificateRequest, error) {
	var req dto.CreateCertificateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateCertificateRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Achievements object Request
func decodeCreateAchievementRequest(r *http.Request) (dto.CreateAchievementRequest, error) {
	var req dto.CreateAchievementRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return dto.CreateAchievementRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

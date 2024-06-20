package handler

import (
	"encoding/json"
	"net/http"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// Decode the

func decodeUserLoginRequest(r *http.Request) (specs.UserLoginRequest, error) {
	var req specs.UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.UserLoginRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Creation object Request
func decodeCreateProfileRequest(r *http.Request) (specs.CreateProfileRequest, error) {
	var req specs.CreateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.CreateProfileRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Updation object Request
func decodeUpdateProfileRequest(r *http.Request) (specs.UpdateProfileRequest, error) {
	var req specs.UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.UpdateProfileRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Education object Request
func decodeCreateEducationRequest(r *http.Request) (specs.CreateEducationRequest, error) {
	var req specs.CreateEducationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.CreateEducationRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Education Updation object Request
func decodeUpdateEducationRequest(r *http.Request) (specs.UpdateEducationRequest, error) {
	var req specs.UpdateEducationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.UpdateEducationRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Project object request
func decodeCreateProjectRequest(r *http.Request) (specs.CreateProjectRequest, error) {
	var req specs.CreateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.CreateProjectRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Projects Updation object Request
func decodeUpdateProjectRequest(r *http.Request) (specs.UpdateProjectRequest, error) {
	var req specs.UpdateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		zap.S().Errorw("error decoding project request", err)
		return specs.UpdateProjectRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Experience object Request
func decodeCreateExperinceRequest(r *http.Request) (specs.CreateExperienceRequest, error) {
	var req specs.CreateExperienceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.CreateExperienceRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Experience Updation object Request
func decodeUpdateExperienceRequest(r *http.Request) (specs.UpdateExperienceRequest, error) {
	var req specs.UpdateExperienceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.UpdateExperienceRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Certicates object Request
func decodeCreateCertificateRequest(r *http.Request) (specs.CreateCertificateRequest, error) {
	var req specs.CreateCertificateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.CreateCertificateRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Certificates Updation object Request
func decodeUpdateCertificateRequest(r *http.Request) (specs.UpdateCertificateRequest, error) {
	var req specs.UpdateCertificateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.UpdateCertificateRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Profile Achievements object Request
func decodeCreateAchievementRequest(r *http.Request) (specs.CreateAchievementRequest, error) {
	var req specs.CreateAchievementRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.CreateAchievementRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

// Decodes the Achievement Updation object Request
func decodeUpdateAchievementRequest(r *http.Request) (specs.UpdateAchievementRequest, error) {
	var req specs.UpdateAchievementRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return specs.UpdateAchievementRequest{}, errors.ErrInvalidBody
	}

	return req, nil
}

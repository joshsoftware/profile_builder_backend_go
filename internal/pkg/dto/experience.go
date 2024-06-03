package dto

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateExperienceRequest struct represents a request to create experience details.
type CreateExperienceRequest struct {
	Experiences []Experience `json:"experiences"`
}

// Experience struct represents details of an experiences.
type Experience struct {
	Designation string `json:"designation"`
	CompanyName string `json:"company_name"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
}

// ExperienceResponse struct represents details of an experiences as response.
type ExperienceResponse struct {
	ID          int    `json:"id"`
	ProfileID   int    `json:"profile_id"`
	Designation string `json:"designation"`
	CompanyName string `json:"company_name"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
}

// ResponseExperience used for response of experiences of profiles
type ResponseExperience struct {
	Experiences []ExperienceResponse `json:"experiences"`
}

// UpdateExperienceRequest represents a request to update education details.
type UpdateExperienceRequest struct {
	Experience Experience `json:"experience"`
}

// Validate func checks if the CreateExperienceRequest is valid.
func (req *CreateExperienceRequest) Validate() error {

	if len(req.Experiences) == 0 {
		return fmt.Errorf("%s : experiences ", errors.ErrEmptyPayload.Error())
	}

	for _, edu := range req.Experiences {
		if edu.Designation == "" {
			return fmt.Errorf("%s : designation ", errors.ErrParameterMissing.Error())
		}

		if edu.CompanyName == "" {
			return fmt.Errorf("%s : company name ", errors.ErrParameterMissing.Error())
		}

		if edu.FromDate == "" {
			return fmt.Errorf("%s : from date ", errors.ErrParameterMissing.Error())
		}

		if edu.ToDate == "" {
			return fmt.Errorf("%s : to date ", errors.ErrParameterMissing.Error())
		}
	}

	return nil
}

// Validate func checks if the UpdateExperienceRequest is valid.
func (req *UpdateExperienceRequest) Validate() error {

	if req.Experience.Designation == "" {
		return fmt.Errorf("%s : designation ", errors.ErrParameterMissing.Error())
	}

	if req.Experience.CompanyName == "" {
		return fmt.Errorf("%s : company name ", errors.ErrParameterMissing.Error())
	}

	if req.Experience.FromDate == "" {
		return fmt.Errorf("%s : from date ", errors.ErrParameterMissing.Error())
	}

	if req.Experience.ToDate == "" {
		return fmt.Errorf("%s : to date ", errors.ErrParameterMissing.Error())
	}

	return nil
}

package dto

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateExperienceRequest struct represents a request to create experience details.
type CreateExperienceRequest struct {
	ProfileID   int          `json:"profile_id"`
	Experiences []Experience `json:"experiences"`
}

// Experience struct represents details of an experiences.
type Experience struct {
	Designation string `json:"designation"`
	CompanyName string `json:"company_name"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
}

// Experience struct represents details of an experiences.
type ExperienceResponse struct {
	ProfileID int `json:"profile_id"`
	Designation string `json:"designation"`
	CompanyName string `json:"company_name"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
}

//ResponseEducation used for response of educations of profiles
type ResponseExperience struct {
	Experiences []ExperienceResponse `json:"experiences"`
}

// Validate func checks if the CreateExperienceRequest is valid.
func (req *CreateExperienceRequest) Validate() error {
	if req.ProfileID == 0 {
		return fmt.Errorf("%s : profile id ", errors.ErrParameterMissing.Error())
	}
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

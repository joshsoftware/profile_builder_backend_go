package specs

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateEducationRequest represents a request to create education details.
type CreateEducationRequest struct {
	Educations []Education `json:"educations"`
}

// UpdateEducationRequest represents a request to update education details.
type UpdateEducationRequest struct {
	Education Education `json:"education"`
}

// ListEducationsFilter used to filter educations based on ids and names
type ListEducationsFilter struct {
	EduationsIDs []int    `json:"educations_ids"`
	Names        []string `json:"names"`
}

// Education represents details of an educational qualification.
type Education struct {
	Degree           string `json:"degree"`
	UniversityName   string `json:"university_name"`
	Place            string `json:"place"`
	PercentageOrCgpa string `json:"percent_or_cgpa"`
	PassingYear      string `json:"passing_year"`
}

// EducationResponse represents details of an educational qualification for specific id.
type EducationResponse struct {
	ID               int    `json:"id"`
	ProfileID        int    `json:"profile_id"`
	Degree           string `json:"degree"`
	UniversityName   string `json:"university_name"`
	Place            string `json:"place"`
	PercentageOrCgpa string `json:"percent_or_cgpa"`
	PassingYear      string `json:"passing_year"`
}

// ResponseEducation used for response of educations of profiles
type ResponseEducation struct {
	Educations []EducationResponse `json:"educations"`
}

// Validate func checks if the CreateEducationRequest is valid.
func (req *CreateEducationRequest) Validate() error {

	if len(req.Educations) == 0 {
		return fmt.Errorf("%s : educations ", errors.ErrEmptyPayload.Error())
	}

	for _, edu := range req.Educations {
		if edu.Degree == "" {
			return fmt.Errorf("%s : degree ", errors.ErrParameterMissing.Error())
		}

		if edu.UniversityName == "" {
			return fmt.Errorf("%s : university name ", errors.ErrParameterMissing.Error())
		}

		if edu.Place == "" {
			return fmt.Errorf("%s : place ", errors.ErrParameterMissing.Error())
		}

		if edu.PercentageOrCgpa == "" {
			return fmt.Errorf("%s : percentage or cgpa ", errors.ErrParameterMissing.Error())
		}

		if edu.PassingYear == "" {
			return fmt.Errorf("%s : passing year ", errors.ErrParameterMissing.Error())
		}
	}

	return nil
}

// Validate func checks if the UpdateEducationRequest is valid.
func (req *UpdateEducationRequest) Validate() error {

	fields := map[string]string{
		"degree":             req.Education.Degree,
		"university name":    req.Education.UniversityName,
		"place":              req.Education.Place,
		"percentage or cgpa": req.Education.PercentageOrCgpa,
		"passing year":       req.Education.PassingYear,
	}

	for fieldName, fieldValue := range fields {
		if fieldValue == "" {
			return fmt.Errorf("%s : %s", errors.ErrParameterMissing.Error(), fieldName)
		}
	}

	return nil
}

type DeleteEducationRequest struct {
	ProfileID   int `json:"profile_id"`
	EducationID int `json:"id"`
}

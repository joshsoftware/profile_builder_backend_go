package dto

import (
	"errors"
	"fmt"
)

type CreateEducationRequest struct {
	ProfileID  int64       `json:"profile_id"`
	Educations []Education `json:"educations"`
}

type Education struct {
	Degree           string `json:"degree"`
	UniversityName   string `json:"university_name"`
	Place            string `json:"place"`
	PercentageOrCgpa string `json:"percent_or_cgpa"`
	PassingYear      string `json:"passing_year"`
}

func (req *CreateEducationRequest) Validate() error {
	if req.ProfileID == 0 {
		return errors.New("profile_id is required")
	}

	for _, edu := range req.Educations {
		if edu.Degree == "" {
			return errors.New("degree is required")
		}

		if edu.UniversityName == "" {
			return errors.New("university name is required")
		}

		if edu.Place == "" {
			return errors.New("place name is required")
		}

		fmt.Println("cgpa:", edu.PercentageOrCgpa)
		if edu.PercentageOrCgpa == "" {
			return errors.New("percentage or cgpa is required")
		}

		if edu.PassingYear == "" {
			return errors.New("passing year is required")
		}
	}

	return nil
}

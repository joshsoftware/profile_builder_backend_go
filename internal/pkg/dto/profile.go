package dto

import (
	"errors"
	"regexp"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
)

type CreateProfileRequest struct {
	Profile Profile `json:"profile"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
type Profile struct {
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	Gender            string   `json:"gender"`
	Mobile            string   `json:"mobile"`
	Designation       string   `json:"designation"`
	Description       string   `json:"description"`
	Title             string   `json:"title"`
	YearsOfExperience float64  `json:"years_of_experience"`
	PrimarySkills     []string `json:"primary_skills"`
	SecondarySkills   []string `json:"secondary_skills"`
	GithubLink        string   `json:"github_link"`
	LinkedinLink      string   `json:"linkedin_link"`
}

func (req *CreateProfileRequest) Validate() error {
	// var errors []string
	if req.Profile.Name == "" {
		return errors.New("name is required")
	}

	if req.Profile.Email == "" {
		return errors.New("email is required")
	}

	matchMail, _ := regexp.MatchString(constants.EMAIL_REGEX, req.Profile.Email)
	if !matchMail {
		return errors.New("invalid email format")
	}

	if req.Profile.Mobile == "" {
		return errors.New("mobile is required")
	}

	matchMob, _ := regexp.MatchString(constants.MOBILE_REGEX, req.Profile.Mobile)
	if !matchMob {
		return errors.New("invalid mobile format")
	}

	if req.Profile.Designation == "" {
		return errors.New("designation is required")
	}

	if req.Profile.Title == "" {
		return errors.New("title is required")
	}

	if req.Profile.YearsOfExperience <= 0.0 {
		return errors.New("years of experience must be a positive number")
	}

	if len(req.Profile.PrimarySkills) == 0 {
		return errors.New("primary skills are required")
	}

	if len(req.Profile.SecondarySkills) == 0 {
		return errors.New("secondary skills are required")
	}

	return nil
}

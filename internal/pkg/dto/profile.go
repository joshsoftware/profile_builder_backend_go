package dto

import (
	"fmt"
	"regexp"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateProfileRequest struct represents a request to create a user profile.
type CreateProfileRequest struct {
	Profile Profile `json:"profile"`
}

// Profile struct represents details of a user profile.
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

// ListProfiles struct represents details of user profiles for listing.
type ListProfiles struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	YearsOfExperience float64  `json:"years_of_experience"`
	PrimarySkills     []string `json:"primary_skills"`
	IsCurrentEmployee int64    `json:"is_current_employee"`
}

// ListProfilesResponse struct represents a response containing a list of user profiles.
type ListProfilesResponse struct {
	Profiles []ListProfiles `json:"profiles"`
}

// Validate func checks if the CreateProfileRequest is valid.
func (req *CreateProfileRequest) Validate() error {

	if req.Profile.Name == "" {
		return fmt.Errorf("%s : name ", errors.ErrParameterMissing.Error())
	}

	if req.Profile.Email == "" {
		return fmt.Errorf("%s : email ", errors.ErrParameterMissing.Error())
	}

	matchMail, _ := regexp.MatchString(constants.EmailRegex, req.Profile.Email)
	if !matchMail {
		return fmt.Errorf("%s : email ", errors.ErrInvalidFormat.Error())
	}

	if req.Profile.Mobile == "" {
		return fmt.Errorf("%s : mobile ", errors.ErrParameterMissing.Error())
	}

	matchMob, _ := regexp.MatchString(constants.MobileRegex, req.Profile.Mobile)
	if !matchMob {
		return fmt.Errorf("%s : mobile ", errors.ErrInvalidFormat.Error())
	}

	if req.Profile.Designation == "" {
		return fmt.Errorf("%s : designation ", errors.ErrParameterMissing.Error())
	}

	if req.Profile.Title == "" {
		return fmt.Errorf("%s : title ", errors.ErrParameterMissing.Error())
	}

	if req.Profile.YearsOfExperience < 0.0 {
		return fmt.Errorf("%s : years of experiences", errors.ErrParameterMissing.Error())
	}

	if len(req.Profile.PrimarySkills) == 0 {
		return fmt.Errorf("%s : primary skills ", errors.ErrParameterMissing.Error())
	}

	if len(req.Profile.SecondarySkills) == 0 {
		return fmt.Errorf("%s : secondary skills ", errors.ErrParameterMissing.Error())
	}

	return nil
}

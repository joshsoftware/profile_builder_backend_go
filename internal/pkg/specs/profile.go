package specs

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

// UpdateProfileRequest struct represents a request to update a user profile.
type UpdateProfileRequest struct {
	ProfileID string  `json:"id"`
	Profile   Profile `json:"profile"`
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
	CareerObjectives  string   `json:"career_objectives"`
}

// ListProfiles struct represents details of user profiles for listing.
type ListProfiles struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	YearsOfExperience float64  `json:"years_of_experience"`
	PrimarySkills     []string `json:"primary_skills"`
	IsCurrentEmployee int      `json:"is_current_employee"`
}

// ResponseListProfiles struct represents response of user profiles for listing.
type ResponseListProfiles struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	YearsOfExperience float64  `json:"years_of_experience"`
	PrimarySkills     []string `json:"primary_skills"`
	IsCurrentEmployee string   `json:"is_current_employee"`
}

// ListSkills struct represents details of skills for listing.
type ListSkills struct {
	Name []string `json:"skills"`
}

// ListProfilesResponse struct represents a response containing a list of user profiles.
type ListProfilesResponse struct {
	Profiles []ResponseListProfiles `json:"profiles"`
}

// ProfileResponse struct represents a response containing profile of specific user.
type ProfileResponse struct {
	Profile ResponseProfile `json:"profile"`
}

// ResponseProfile struct represents details of a user profile as in response.
type ResponseProfile struct {
	ProfileID         int      `json:"id"`
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
	CareerObjectives  string   `json:"career_objectives"`
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
	if req.Profile.CareerObjectives == "" {
		return fmt.Errorf("%s : career objectives ", errors.ErrParameterMissing.Error())
	}

	return nil
}

// Validate func checks if the UpdateProfileRequest is valid.
func (req *UpdateProfileRequest) Validate() error {

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

package specs

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	JoshJoiningDate   string   `json:"josh_joining_date"`
	GithubLink        string   `json:"github_link"`
	LinkedinLink      string   `json:"linkedin_link"`
	CareerObjectives  string   `json:"career_objectives"`
}

// ListProfiles struct represents details of user profiles for listing.
type ListProfiles struct {
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	YearsOfExperience float64        `json:"years_of_experience"`
	PrimarySkills     []string       `json:"primary_skills"`
	IsCurrentEmployee int            `json:"is_current_employee"`
	IsActive          int            `json:"is_active"`
	JoshJoiningDate   sql.NullString `json:"josh_joining_date"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	IsProfileComplete int            `json:"is_profile_complete"`
}

// ResponseListProfiles struct represents response of user profiles for listing.
type ResponseListProfiles struct {
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	YearsOfExperience float64        `json:"years_of_experience"`
	PrimarySkills     []string       `json:"primary_skills"`
	IsCurrentEmployee string         `json:"is_current_employee"`
	IsActive          string         `json:"is_active"`
	JoshJoiningDate   sql.NullString `json:"josh_joining_date"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	IsProfileComplete string         `json:"is_profile_complete"`
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
	ProfileID         int            `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	Gender            string         `json:"gender"`
	Mobile            string         `json:"mobile"`
	Designation       string         `json:"designation"`
	Description       string         `json:"description"`
	Title             string         `json:"title"`
	YearsOfExperience float64        `json:"years_of_experience"`
	PrimarySkills     []string       `json:"primary_skills"`
	SecondarySkills   []string       `json:"secondary_skills"`
	JoshJoiningDate   sql.NullString `json:"josh_joining_date"`
	GithubLink        string         `json:"github_link"`
	LinkedinLink      string         `json:"linkedin_link"`
	CareerObjectives  string         `json:"career_objectives"`
	IsInvited         string         `json:"is_invited"`
}

// UpdateSequenceRequest struct represents a request to update a sequence of component.
type UpdateSequenceRequest struct {
	ProfileID int `json:"id"`
	Component `json:"component"`
}

// Component used while update request of sequence of component
type Component struct {
	CompName            string      `json:"comp_name"`
	ComponentPriorities map[int]int `json:"component_priorities"`
}

// UpdateProfileStatus struct represents a request to update a status of profile.
type UpdateProfileStatus struct {
	ProfileStatus UpdateProfileStatusRequest `json:"profile_status"`
}

// UpdateProfileStatusRequest struct represents a request to update a status of profile.
type UpdateProfileStatusRequest struct {
	IsCurrentEmployee string `json:"is_current_employee,omitempty"`
	IsActive          string `json:"is_active,omitempty"`
}

// UserContext struct represents a user context.
type UserContext struct {
	Role  string `json:"role"`
	Email string `json:"email"`
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

	if req.Profile.Title == "" {
		return fmt.Errorf("%s : title ", errors.ErrParameterMissing.Error())
	}

	if req.Profile.YearsOfExperience < 0.0 {
		return fmt.Errorf("%s : years of experiences", errors.ErrParameterMissing.Error())
	}

	if req.Profile.Description == "" {
		return fmt.Errorf("%s : description ", errors.ErrParameterMissing.Error())
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

	if req.Profile.Title == "" {
		return fmt.Errorf("%s : title ", errors.ErrParameterMissing.Error())
	}

	if req.Profile.YearsOfExperience < 0.0 {
		return fmt.Errorf("%s : years of experiences", errors.ErrParameterMissing.Error())
	}

	if req.Profile.Description == "" {
		return fmt.Errorf("%s : description ", errors.ErrParameterMissing.Error())
	}

	return nil
}

// Validate func checks if the UpdateSequenceRequest is valid.
func (req *UpdateSequenceRequest) Validate() error {

	if req.ProfileID <= 0 {
		return fmt.Errorf("%s : profile id", errors.ErrParameterMissing.Error())
	}

	if req.Component.CompName == "" {
		return fmt.Errorf("%s : component name", errors.ErrParameterMissing.Error())
	}

	if !constants.ComponentMap[strings.ToLower(req.Component.CompName)] {
		return fmt.Errorf("%s : %s", errors.ErrComponentNotSuppoerted.Error(), req.Component.CompName)
	}

	if len(req.Component.ComponentPriorities) == 0 {
		return fmt.Errorf("%s : component priorities", errors.ErrParameterMissing.Error())
	}

	for key, value := range req.Component.ComponentPriorities {
		if key < 0 || value < 0 {
			return fmt.Errorf("%s : component priorities (keys and values must be non-negative)", errors.ErrInvalidFormat.Error())
		}
	}

	return nil
}

// Validate func checks if the UpdateProfileStatusRequest is valid.
func (req *UpdateProfileStatusRequest) Validate() error {
	if req.IsCurrentEmployee == "" && req.IsActive == "" {
		return fmt.Errorf("%s : at least one of is_current_employee or is_active must be provided", errors.ErrParameterMissing.Error())
	}

	if req.IsCurrentEmployee != "" {
		normalizedValue := strings.ToUpper(req.IsCurrentEmployee)
		if !(normalizedValue == "YES" || normalizedValue == "NO") {
			return fmt.Errorf("%s : is_current_employee must be 'YES' or 'NO'", errors.ErrInvalidBody.Error())
		}
	}

	if req.IsActive != "" {
		normalizedValue := strings.ToUpper(req.IsActive)
		if !(normalizedValue == "YES" || normalizedValue == "NO") {
			return fmt.Errorf("%s : is_active must be 'YES' or 'NO'", errors.ErrInvalidBody.Error())
		}
	}
	return nil
}

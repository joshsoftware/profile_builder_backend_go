package specs

import (
	"fmt"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateFullProfileRequest represents a request to create a profile along with educations and projects.
type CreateFullProfileRequest struct {
	Profile    Profile     `json:"profile"`
	Educations []Education `json:"educations,omitempty"`
	Projects   []Project   `json:"projects,omitempty"`
}

// Validate checks if the CreateFullProfileRequest is valid.
// Relaxed validation is used for educations and projects.
func (req *CreateFullProfileRequest) Validate() error {
	profileReq := CreateProfileRequest{Profile: req.Profile}
	if err := profileReq.Validate(); err != nil {
		return err
	}

	for _, edu := range req.Educations {
		if edu.Degree == "" {
			return fmt.Errorf("%s : degree ", errors.ErrParameterMissing.Error())
		}
	}

	for _, proj := range req.Projects {
		if proj.Name == "" {
			return fmt.Errorf("%s : project name ", errors.ErrParameterMissing.Error())
		}
	}

	return nil
}

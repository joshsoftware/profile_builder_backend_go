package dto

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateAchievementRequest struct represents a request to create achievements details.
type CreateAchievementRequest struct {
	Achievements []Achievement `json:"achievements"`
}

// Achievement struct represents details of an achievements.
type Achievement struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateAchievementRequest struct represents a request to update a achievement
type UpdateAchievementRequest struct {
	Achievement Achievement `json:"achievement"`
}

// Validate func checks if the CreateAchievementRequest is valid.
func (req *CreateAchievementRequest) Validate() error {

	if len(req.Achievements) == 0 {
		return fmt.Errorf("%s : achievements ", errors.ErrEmptyPayload.Error())
	}

	for _, edu := range req.Achievements {
		if edu.Name == "" {
			return fmt.Errorf("%s : name ", errors.ErrParameterMissing.Error())
		}

		if edu.Description == "" {
			return fmt.Errorf("%s : decsription ", errors.ErrParameterMissing.Error())
		}
	}

	return nil
}

// Validate func checks if the UpdateAchievementRequest is valid.
func (req *UpdateAchievementRequest) Validate() error {

	if req.Achievement.Name == "" {
		return fmt.Errorf("%s : name ", errors.ErrParameterMissing.Error())
	}

	if req.Achievement.Description == "" {
		return fmt.Errorf("%s : decsription ", errors.ErrParameterMissing.Error())
	}

	return nil
}

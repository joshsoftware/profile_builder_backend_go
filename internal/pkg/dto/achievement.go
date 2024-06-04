package dto

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateAchievementRequest struct represents a request to create achievements details.
type CreateAchievementRequest struct {
	ProfileID    int           `json:"profile_id"`
	Achievements []Achievement `json:"achievements"`
}

// Achievement struct represents details of an achievements.
type Achievement struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AchievementResponse struct {
	ProfileID   int    `json:"profile_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseAchievement struct {
	Achievements []AchievementResponse `json:"achievements"`
}

// Validate func checks if the CreateAchievementRequest is valid.
func (req *CreateAchievementRequest) Validate() error {
	if req.ProfileID == 0 {
		return fmt.Errorf("%s : profile id ", errors.ErrParameterMissing.Error())
	}
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

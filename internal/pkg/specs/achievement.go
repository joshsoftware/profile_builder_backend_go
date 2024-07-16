package specs

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// ListAchievementFilter used to filter achievements based on ids and names
type ListAchievementFilter struct {
	AchievementIDs []int    `json:"achievement_ids"`
	Names          []string `json:"names"`
}

// CreateAchievementRequest struct represents a request to create achievements details.
type CreateAchievementRequest struct {
	Achievements []Achievement `json:"achievements"`
}

// UpdateAchievementRequest struct represents a request to update a achievement
type UpdateAchievementRequest struct {
	Achievement Achievement `json:"achievement"`
}

// Achievement struct represents details of an achievements.
type Achievement struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AchievementResponse struct represents details of achievement response
type AchievementResponse struct {
	ID          int    `json:"id"`
	ProfileID   int    `json:"profile_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ResponseAchievement struct represents array of achievements which should be returned
type ResponseAchievement struct {
	Achievements []AchievementResponse `json:"achievements"`
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
	}

	return nil
}

// Validate func checks if the UpdateAchievementRequest is valid.
func (req *UpdateAchievementRequest) Validate() error {

	if req.Achievement.Name == "" {
		return fmt.Errorf("%s : name ", errors.ErrParameterMissing.Error())
	}
	return nil
}

type DeleteAchievementRequest struct {
	ProfileID     int `json:"profile_id"`
	AchievementID int `json:"id"`
}

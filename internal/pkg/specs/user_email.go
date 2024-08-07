package specs

import "time"

type InvitationResponse struct {
	ProfileID       int       `json:"profile_id"`
	ProfileComplete int       `json:"is_profile_complete"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedByID     int       `json:"created_by_id"`
	UpdatedByID     int       `json:"updated_by_id"`
}

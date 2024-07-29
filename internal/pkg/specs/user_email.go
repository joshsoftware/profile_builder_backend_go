package specs

import (
	"fmt"

	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// UserEmailRequest is the request for sending an email to the user
type UserSendInvitationRequest struct {
	ProfileID int `json:"profile_id"`
}

// Validate validates the request
func (req *UserSendInvitationRequest) Validate() error {
	if req.ProfileID <= 0 {
		return fmt.Errorf("%s : profile id", errors.ErrParameterMissing.Error())
	}
	return nil
}

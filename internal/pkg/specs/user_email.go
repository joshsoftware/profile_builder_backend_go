package specs

import (
	"fmt"
	"regexp"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

type UserEmailRequest struct {
	ProfileID int    `json:"profile_id"`
	Email     string `json:"user_email"`
}

func (req *UserEmailRequest) Validate() error {
	if req.ProfileID <= 0 {
		return fmt.Errorf("%s : profile id", errors.ErrParameterMissing.Error())
	}

	if req.Email == "" {
		return fmt.Errorf("%s : email ", errors.ErrParameterMissing.Error())
	}

	matchMail, _ := regexp.MatchString(constants.EmailRegex, req.Email)
	if !matchMail {
		return fmt.Errorf("%s : email ", errors.ErrInvalidFormat.Error())
	}

	return nil
}

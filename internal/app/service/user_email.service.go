package service

import (
	"context"
	"net/smtp"
	"os"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

type UserEmailService interface {
	SendUserInvitation(ctx context.Context, request specs.UserEmailRequest) error
}

func getEmailConfig() (from, password, smtpServer, smtpPort string) {
	from = os.Getenv("FROM_EMAIL")
	password = os.Getenv("FROM_EMAIL_PASSWORD")
	smtpServer = os.Getenv("EMAIL_HOST")
	smtpPort = os.Getenv("EMAIL_HOST_PORT")
	return
}

func (userService *service) SendUserInvitation(ctx context.Context, request specs.UserEmailRequest) error {
	from, password, smtpServer, smtpPort := getEmailConfig()

	// Setup the authentication
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// Construct the email message
	message := helpers.ConstructEmailMessage(from, request)

	// Send the email
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{request.Email}, message)
	if err != nil {
		zap.S().Error("Error sending mail: ", err)
		return errors.ErrUnableToSendEmail
	}
	zap.S().Info("Email sent successfully to: ", request.Email, " with profile: ", request.ProfileID)
	return nil
}

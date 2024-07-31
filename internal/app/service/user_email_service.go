package service

import (
	"context"
	"time"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// UserEmailService is the interface for the user email service
type UserEmailService interface {
	SendUserInvitation(ctx context.Context, userID int, request specs.UserSendInvitationRequest) error
	SendAdminInvitation(ctx context.Context, userID int, request specs.UserSendInvitationRequest) error
}

// inviteUser sends an email to the user with the invitation link
func (userService *service) SendUserInvitation(ctx context.Context, userID int, request specs.UserSendInvitationRequest) (err error) {
	tx, _ := userService.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := userService.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	sendRequest := repository.SendUserInvitationRequest{
		ProfileID: request.ProfileID,
	}
	email, err := userService.UserEmailRepo.GetEmailByProfileID(ctx, sendRequest, tx)
	if err != nil {
		zap.S().Error("Error getting email by profile id: ", err)
		return err
	}

	now := helpers.GetCurrentISTTime()
	createSendInvitationRequest := repository.EmailRepo{
		ProfileID:       request.ProfileID,
		ProfileComplete: 0,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedByID:     userID,
		UpdatedByID:     userID,
	}

	err = userService.UserEmailRepo.CreateSendInvitation(ctx, createSendInvitationRequest, tx)
	if err != nil {
		zap.S().Error("Error creating send invitation: ", err)
		return err
	}

	err = userService.UserLoginRepo.CreateUserAsEmployee(ctx, email, tx)
	if err != nil {
		zap.S().Error("Error creating user as employee: ", err)
		return err
	}

	go func(email string, profileID int) {
		for i := 0; i < constants.DefaultMaxRetries; i++ {
			err := helpers.SendUserInvitation(email, profileID)
			if err != nil {
				zap.S().Error("Error sending invitation to user : ", err)
				time.Sleep(time.Second * 2)
				continue
			}
			return
		}
		zap.S().Error("Max retries reached. Failed to send invitation.")
	}(email, request.ProfileID)

	return nil
}

func (userService *service) SendAdminInvitation(ctx context.Context, userID int, request specs.UserSendInvitationRequest) (err error) {
	tx, _ := userService.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := userService.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	sendRequest := repository.SendUserInvitationRequest{
		ProfileID: request.ProfileID,
	}

	ID, err := userService.UserEmailRepo.GetCreatedByIdByProfileID(ctx, sendRequest, tx)
	if err != nil {
		zap.S().Error("Error getting created_by_id by profile id: ", err)
		return err
	}

	adminEmail, err := userService.UserEmailRepo.GetUserEmailByUserID(ctx, ID, tx)
	if err != nil {
		zap.S().Error("Error getting email by id: ", err)
		return err
	}
	now := helpers.GetCurrentISTTime()
	updateSendRequest := repository.UpadateRequest{
		ProfileID: request.ProfileID,
		UpdatedAt: now,
	}
	err = userService.UserEmailRepo.UpdateProfileCompleteStatus(ctx, updateSendRequest, tx)
	if err != nil {
		zap.S().Error("Error creating send invitation: ", err)
		return err
	}

	employeeEmail, err := userService.UserEmailRepo.GetEmailByProfileID(ctx, sendRequest, tx)
	if err != nil {
		zap.S().Error("Error getting email by profile id: ", err)
		return err
	}

	err = userService.UserLoginRepo.RemoveUserEmployee(ctx, employeeEmail, tx)
	if err != nil {
		zap.S().Error("Error removing user employee: ", err)
		return err
	}

	// Send email to the user in a fire-and-forget manner
	go func(adminEmail string, profileID int) {
		for i := 0; i < constants.DefaultMaxRetries; i++ {
			err := helpers.SendAdminInvitation(adminEmail, profileID)
			if err != nil {
				zap.S().Error("Error sending invitation: ", err)
				time.Sleep(time.Second * 2)
				continue
			}
			return
		}
		zap.S().Error("Max retries reached. Failed to send invitation.")
	}(adminEmail, request.ProfileID)

	return nil
}

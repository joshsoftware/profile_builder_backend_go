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

	createSendInvitationRequest := repository.EmailRepo{
		ProfileID:       request.ProfileID,
		ProfileComplete: 0,
		CreatedAt:       today,
		UpdatedAt:       today,
		CreatedByID:     userID,
		UpdatedByID:     userID,
	}

	err = userService.UserEmailRepo.CreateSendInvitation(ctx, createSendInvitationRequest, tx)
	if err != nil {
		zap.S().Error("Error creating send invitation: ", err)
		return err
	}

	go func(email string, profileID int) {
		for i := 0; i < constants.DefaultMaxRetries; i++ {
			helpers.SendUserInvitation(email, profileID)
			if err != nil {
				zap.S().Error("Error sending invitation: ", err)
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

	email, err := userService.UserEmailRepo.GetUserEmailByUserID(ctx, ID, tx)
	if err != nil {
		zap.S().Error("Error getting email by id: ", err)
		return err
	}

	updateSendRequest := repository.SendUserInvitationRequest{
		ProfileID: request.ProfileID,
	}
	err = userService.UserEmailRepo.UpdateProfileCompleteStatus(ctx, updateSendRequest, tx)
	if err != nil {
		zap.S().Error("Error creating send invitation: ", err)
		return err
	}

	// Send email to the user in a fire-and-forget manner
	go func(email string) {
		for i := 0; i < constants.DefaultMaxRetries; i++ {
			helpers.SendAdminInvitation(email, request.ProfileID)
			if err != nil {
				zap.S().Error("Error sending invitation: ", err)
				time.Sleep(time.Second * 2)
				continue
			}
			return
		}
		zap.S().Error("Max retries reached. Failed to send invitation.")
	}(email)

	return nil
}

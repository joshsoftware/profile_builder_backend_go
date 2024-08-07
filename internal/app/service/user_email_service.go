package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// UserEmailService is the interface for the user email service
type UserEmailService interface {
	SendUserInvitation(ctx context.Context, userID int, profileID int) error
	UpdateInvitation(ctx context.Context, userID int, profileID int) error
}

// inviteUser sends an email to the user with the invitation link
func (userService *service) SendUserInvitation(ctx context.Context, userID int, profileID int) (err error) {
	tx, _ := userService.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := userService.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	profile, err := userService.ProfileRepo.GetProfile(ctx, profileID, tx)
	if err != nil {
		zap.S().Error("Error getting profile by profile id: ", err)
		return err
	}

	err = helpers.SendUserInvitation(profile.Email, profileID)
	if err != nil {
		zap.S().Error("Error sending invitation to user : ", err)
		return err
	}

	now := helpers.GetCurrentISTTime()
	createInvitationRequest := repository.Invitations{
		ProfileID:       profileID,
		ProfileComplete: constants.ProfileIncomplete,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedByID:     userID,
		UpdatedByID:     userID,
	}

	err = userService.UserEmailRepo.CreateInvitation(ctx, createInvitationRequest, tx)
	if err != nil {
		zap.S().Error("Error creating send invitation: ", err)
		return err
	}

	err = userService.UserLoginRepo.CreateUser(ctx, profile.Email, constants.Employee, tx)
	if err != nil {
		zap.S().Error("Error creating user as employee: ", err)
		return err
	}

	zap.S().Info("Invitation sent to user")
	return nil
}

// Update profile complete status with sending the email to the admin
func (userService *service) UpdateInvitation(ctx context.Context, userID int, profileID int) (err error) {
	tx, _ := userService.ProfileRepo.BeginTransaction(ctx)
	defer func() {
		txErr := userService.ProfileRepo.HandleTransaction(ctx, tx, err)
		if txErr != nil {
			err = txErr
			return
		}
	}()

	invitation, err := userService.UserEmailRepo.GetInvitations(ctx, profileID, tx)
	if err != nil {
		zap.S().Error("Error getting created_by_id by profile id: ", err)
		return err
	}

	userInfoFilter := specs.UserInfoFilter{
		ID: invitation.CreatedByID,
	}

	admin, err := userService.UserLoginRepo.GetUserInfo(ctx, userInfoFilter)
	if err != nil {
		zap.S().Error("Error getting email by id: ", err)
		return err
	}

	err = helpers.SendAdminInvitation(admin.Email, profileID)
	if err != nil {
		zap.S().Error("Error sending invitation: ", err)
		return err
	}

	now := helpers.GetCurrentISTTime()
	updateSendRequest := repository.UpadateRequest{
		ProfileComplete: constants.ProfileComplete,
		UpdatedAt:       now,
	}
	err = userService.UserEmailRepo.UpdateProfileCompleteStatus(ctx, profileID, updateSendRequest, tx)
	if err != nil {
		zap.S().Error("Error creating send invitation: ", err)
		return err
	}

	profile, err := userService.ProfileRepo.GetProfile(ctx, profileID, tx)
	if err != nil {
		zap.S().Error("Error getting email by profile id: ", err)
		return err
	}

	err = userService.UserLoginRepo.RemoveUser(ctx, profile.Email, tx)
	if err != nil {
		zap.S().Error("Error removing user employee: ", err)
		return err
	}

	zap.S().Info("Profile completed successfully")
	return nil
}

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
		zap.S().Errorf("Error while getting profile. profile id :%d, err:%v", profileID, err)
		return err
	}

	err = helpers.SendUserInvitation(profile.Email, profileID)
	if err != nil {
		zap.S().Errorf("Error sending invitation:%v  for email:%s and profile ID : %d ", err, profile.Email, profileID)
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
		zap.S().Errorf("Error creating send invitation %v. for user %s : ", err, profile.Email)
		return err
	}

	err = userService.UserLoginRepo.CreateUser(ctx, profile.Email, constants.Employee, tx)
	if err != nil {
		zap.S().Errorf("Error creating user employee: %v for user %s: ", err, profile.Email)
		return err
	}

	zap.S().Infof("Invitation sent to user : user ID : %d and profile ID : %d", userID, profileID)
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

	getRequest := repository.GetRequest{
		ProfileID:         profileID,
		IsProfileComplete: constants.ProfileIncomplete,
	}

	invitation, err := userService.UserEmailRepo.GetInvitations(ctx, getRequest, tx)
	if err != nil {
		zap.S().Errorf("Error getting invitation : %v by profile ID: %d", err, profileID)
		return err
	}

	userInfoFilter := specs.UserInfoFilter{
		ID: invitation.CreatedByID,
	}

	admin, err := userService.UserLoginRepo.GetUserInfo(ctx, userInfoFilter)
	if err != nil {
		zap.S().Errorf("Error getting email : %v by user : %s: ", err, invitation.CreatedByID)
		return err
	}

	err = helpers.SendAdminInvitation(admin.Email, profileID)
	if err != nil {
		zap.S().Errorf("Error sending invitation %v for user email : %s and profile ID : %d: ", err, admin.Email, profileID)
		return err
	}

	now := helpers.GetCurrentISTTime()
	updateSendRequest := repository.UpdateRequest{
		ProfileComplete: constants.ProfileComplete,
		UpdatedAt:       now,
	}
	err = userService.UserEmailRepo.UpdateProfileCompleteStatus(ctx, profileID, updateSendRequest, tx)
	if err != nil {
		zap.S().Errorf("Error creating send invitation %v for profile ID : %d : ", err, profileID)
		return err
	}

	profile, err := userService.ProfileRepo.GetProfile(ctx, profileID, tx)
	if err != nil {
		zap.S().Errorf("Error getting email : %v by profile id: %d ", err, profileID)
		return err
	}

	err = userService.UserLoginRepo.RemoveUser(ctx, profile.Email, tx)
	if err != nil {
		zap.S().Errorf("Error removing user employee: %v for user : %s ", err, profile.Email)
		return err
	}

	zap.S().Infof("Profile completed successfully for user : user ID : %d and profile ID : %d", userID, profileID)
	return nil
}

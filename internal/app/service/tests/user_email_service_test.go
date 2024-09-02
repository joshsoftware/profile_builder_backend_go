package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/undefinedlabs/go-mpatch"
)

var mockTx pgx.Tx

type ServiceTestSuite struct {
	suite.Suite
	emailService service.Service
	emailRepo    *mocks.EmailStorer
	loginRepo    *mocks.UserStorer
	profileRepo  *mocks.ProfileStorer
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) SetupTest() {
	s.emailRepo = &mocks.EmailStorer{}
	s.profileRepo = &mocks.ProfileStorer{}
	s.loginRepo = &mocks.UserStorer{}
	s.emailService = service.NewServices(service.RepoDeps{
		UserEmailDeps: s.emailRepo,
		UserLoginDeps: s.loginRepo,
		ProfileDeps:   s.profileRepo,
	})
}

func (s *ServiceTestSuite) TearDownSuite() {
	s.emailRepo.AssertExpectations(s.T())
	s.profileRepo.AssertExpectations(s.T())
	s.loginRepo.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestSendUserInvitation() {

	mpatch.PatchMethod(helpers.GetCurrentISTTime, func() string {
		return "2024-08-20T18:42:34+05:30"
	})

	type args struct {
		ctx       context.Context
		profileID int
		err       error
	}

	var (
		UserID    = 1
		ProfileID = 1
	)

	mockResponseProfile := specs.ResponseProfile{
		ProfileID:         ProfileID,
		Name:              "Test User",
		Email:             "test@example.com",
		Gender:            "Male",
		Mobile:            "1234567890",
		Designation:       "Software Engineer",
		Description:       "Test Description",
		Title:             "Test Title",
		YearsOfExperience: 2.5,
		PrimarySkills:     []string{"Go", "Python"},
		SecondarySkills:   []string{"React", "Angular"},
		JoshJoiningDate:   sql.NullString{String: "2021-01-01", Valid: true},
		GithubLink:        "github.com/test",
		LinkedinLink:      "linkedin.com/test",
		CareerObjectives:  "Test Career Objectives",
	}

	mockInvitationRequest := repository.Invitations{
		ProfileID:       mockResponseProfile.ProfileID,
		ProfileComplete: constants.ProfileIncomplete,
		CreatedAt:       helpers.GetCurrentISTTime(),
		UpdatedAt:       helpers.GetCurrentISTTime(),
		CreatedByID:     UserID,
		UpdatedByID:     UserID,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args args)
	}{
		// POSITIVE || successfully send invitation
		{
			name: "Successfully send invitation",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       nil,
			},
			wantErr: false,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, nil).Once()

				patch, _ := mpatch.PatchMethod(helpers.SendUserInvitation, func(email, name string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()

				s.emailRepo.On("CreateInvitation", args.ctx, mockInvitationRequest, mock.Anything).Return(nil).Once()
				s.loginRepo.On("CreateUser", args.ctx, mockResponseProfile.Name, mockResponseProfile.Email, constants.Employee, mock.Anything).Return(nil).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, nil).Return(nil).Once()
			},
		},
		// NEGATIVE || failed to get profile
		{
			name: "Failed to get profile",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to get profile"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},

		// NEGATIVE || failed to create invitation
		{
			name: "Failed to create invitation",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to create invitation"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, nil).Once()

				patch, _ := mpatch.PatchMethod(helpers.SendUserInvitation, func(email, name string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()

				s.emailRepo.On("CreateInvitation", args.ctx, mockInvitationRequest, mock.Anything).Return(args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},

		// NEGATIVE || failed to create user
		{
			name: "Failed to create user",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to create user"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, nil).Once()

				patch, _ := mpatch.PatchMethod(helpers.SendUserInvitation, func(email, name string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()

				s.emailRepo.On("CreateInvitation", args.ctx, mockInvitationRequest, mock.Anything).Return(nil).Once()
				s.loginRepo.On("CreateUser", args.ctx, mockResponseProfile.Name, mockResponseProfile.Email, constants.Employee, mock.Anything).Return(args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.prepare(tt.args)
			err := s.emailService.SendUserInvitation(tt.args.ctx, UserID, ProfileID)
			assert.Equal(s.T(), tt.wantErr, err != nil)
		})
	}
}

func (s *ServiceTestSuite) TestUpdateInvitation() {
	mpatch.PatchMethod(time.Now(), func() string {
		return " 2009-11-10 23:00:00 +0000 UTC m=+0.000000001"
	})

	mpatch.PatchMethod(helpers.GetCurrentISTTime, func() string {
		return "2024-08-20T18:42:34+05:30"
	})

	type args struct {
		ctx       context.Context
		profileID int
		err       error
	}

	var (
		UserID    = 1
		ProfileID = 1
	)
	mockUserInfoFilter := specs.UserInfoFilter{
		ID: UserID,
	}

	mockRequest := repository.GetRequest{
		ProfileID:         ProfileID,
		IsProfileComplete: constants.ProfileIncomplete,
	}

	mockUpdateRequest := repository.UpdateRequest{
		ProfileComplete: constants.ProfileComplete,
		UpdatedAt:       helpers.GetCurrentISTTime(),
	}
	mockAdminInfo := repository.User{
		ID:    int64(profileID),
		Email: "admin@example.com",
		Role:  constants.Admin,
	}
	mockInvitationRequest := specs.InvitationResponse{
		ProfileID:       mockResponseProfile.ProfileID,
		ProfileComplete: constants.ProfileComplete,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CreatedByID:     UserID,
		UpdatedByID:     UserID,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args args)
	}{
		// POSITIVE || successfully update invitation
		{
			name: "Successfully update invitation",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       nil,
			},
			wantErr: false,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.emailRepo.On("GetInvitations", args.ctx, mockRequest, mock.Anything).Return(mockInvitationRequest, nil).Once()
				s.loginRepo.On(("GetUserInfo"), args.ctx, mockUserInfoFilter).Return(mockAdminInfo, nil).Once()
				patch, _ := mpatch.PatchMethod(helpers.SendAdminInvitation, func(email, name string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()
				s.emailRepo.On("UpdateProfileCompleteStatus", args.ctx, args.profileID, mockUpdateRequest, mock.Anything).Return(nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, nil).Once()
				s.loginRepo.On("RemoveUser", args.ctx, mockResponseProfile.Email, mock.Anything).Return(nil).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, nil).Return(nil).Once()
			},
		},

		// NEGATIVE || failed to get invitation
		{
			name: "Failed to get invitation",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to get invitation"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.emailRepo.On("GetInvitations", args.ctx, mockRequest, mock.Anything).Return(mockInvitationRequest, args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},
		{
			name: "Failed to get admin",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to get admin"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.emailRepo.On("GetInvitations", args.ctx, mockRequest, mock.Anything).Return(mockInvitationRequest, nil).Once()
				s.loginRepo.On(("GetUserInfo"), args.ctx, mockUserInfoFilter).Return(mockAdminInfo, args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},

		// NEGATIVE || failed to update profile complete status
		{
			name: "Failed to update profile complete status",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to update profile complete status"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.emailRepo.On("GetInvitations", args.ctx, mockRequest, mock.Anything).Return(mockInvitationRequest, nil).Once()
				s.loginRepo.On(("GetUserInfo"), args.ctx, mockUserInfoFilter).Return(mockAdminInfo, nil).Once()
				patch, _ := mpatch.PatchMethod(helpers.SendAdminInvitation, func(email string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()
				s.emailRepo.On("UpdateProfileCompleteStatus", args.ctx, args.profileID, mockUpdateRequest, mock.Anything).Return(args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},

		// NEGATIVE || failed to get profile
		{
			name: "Failed to get profile",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to get profile"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.emailRepo.On("GetInvitations", args.ctx, mockRequest, mock.Anything).Return(mockInvitationRequest, nil).Once()
				s.loginRepo.On(("GetUserInfo"), args.ctx, mockUserInfoFilter).Return(mockAdminInfo, nil).Once()
				patch, _ := mpatch.PatchMethod(helpers.SendAdminInvitation, func(email, name string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()
				s.emailRepo.On("UpdateProfileCompleteStatus", args.ctx, args.profileID, mockUpdateRequest, mock.Anything).Return(nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},

		// NEGATIVE || failed to remove user
		{
			name: "Failed to remove user",
			args: args{
				ctx:       context.Background(),
				profileID: ProfileID,
				err:       errors.New("failed to remove user"),
			},
			wantErr: true,
			prepare: func(args args) {
				s.profileRepo.On("BeginTransaction", args.ctx).Return(mockTx, nil).Once()
				s.emailRepo.On("GetInvitations", args.ctx, mockRequest, mock.Anything).Return(mockInvitationRequest, nil).Once()
				s.loginRepo.On(("GetUserInfo"), args.ctx, mockUserInfoFilter).Return(mockAdminInfo, nil).Once()
				patch, _ := mpatch.PatchMethod(helpers.SendAdminInvitation, func(email, name string, profileID int) error {
					return nil
				})
				defer patch.Unpatch()
				s.emailRepo.On("UpdateProfileCompleteStatus", args.ctx, args.profileID, mockUpdateRequest, mock.Anything).Return(nil).Once()
				s.profileRepo.On("GetProfile", args.ctx, args.profileID, mock.Anything).Return(mockResponseProfile, nil).Once()
				s.loginRepo.On("RemoveUser", args.ctx, mockResponseProfile.Email, mock.Anything).Return(args.err).Once()
				s.profileRepo.On("HandleTransaction", args.ctx, mock.Anything, args.err).Return(args.err).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.prepare(tt.args)
			err := s.emailService.UpdateInvitation(tt.args.ctx, UserID, ProfileID)
			assert.Equal(s.T(), tt.wantErr, err != nil)
		})
	}

}

package profile

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestGenerateLoginToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "SECRET_KEY")
	defer os.Unsetenv("SECRET_KEY")

	mockProfileRepo := mocks.NewProfileStorer(t)
	profileService := NewServices(mockProfileRepo)
	tests := []struct {
		name    string
		email   string
		setup   func(userMock *mocks.ProfileStorer)
		isError bool
	}{
		{
			name:  "success",
			email: "saurabh.puri@joshsoftware.com",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("GetProfileByEmail", mock.Anything, "saurabh.puri@joshsoftware.com").Return(int64(1), nil)
			},
			isError: false,
		},
		{
			name:  "failure",
			email: "saurabh.puri55@joshsoftware.com",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("GetProfileByEmail", mock.Anything, "saurabh.puri55@joshsoftware.com").Return(int64(0), errors.New("error"))
			},
			isError: true,
		},
		{
			name:  "failure user not found",
			email: "saurabh.puri55@joshsoftware.com",
			setup: func(userMock *mocks.ProfileStorer) {
				userMock.On("GetProfileByEmail", mock.Anything, "saurabh.puri55@joshsoftware.com").Return(int64(0), errors.New("user not found"))
			},
			isError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(mockProfileRepo)
			_, err := profileService.GenerateLoginToken(context.TODO(), test.email)
			if (err != nil) != test.isError {
				t.Errorf("Test Failed , expected error %v but got %v", test.isError, err)
			}
		})
	}
}

func Test_service_CreateToken(t *testing.T) {
	type fields struct {
		Repo repository.ProfileStorer
	}

	type args struct {
		email  string
		userId int64
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid inputs",
			fields: fields{
				Repo: nil,
			},
			args: args{
				email:  "saurabh.puri@joshsoftware.com",
				userId: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid secret key",
			fields: fields{
				Repo: nil,
			},
			args: args{
				email:  "saurabh.puri@joshsoftware.com",
				userId: 1,
			},
			wantErr: true,
		},
		{
			name: "failed to create token",
			fields: fields{
				Repo: nil,
			},
			args: args{
				email:  "saurabh.puri@joshsoftware.com",
				userId: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profileSvc := &service{
				Repo: tt.fields.Repo,
			}
			if tt.name == "invalid secret key" {
				os.Setenv("SECRET_KEY", "")
			} else if tt.name == "valid inputs" {
				if os.Getenv("SECRET_KEY") == "" {
					os.Setenv("SECRET_KEY", "SECRET_KEY")
				}
			} else if tt.name == "failed to create token" {
				os.Setenv("SECRET_KEY", "")
			}

			_, err := profileSvc.createToken(tt.args.email, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.createToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	"go.uber.org/zap"
)

// service implements the Service interface.
type service struct {
	UserLoginRepo   repository.UserStorer
	ProfileRepo     repository.ProfileStorer
	EducationRepo   repository.EducationStorer
	ExperienceRepo  repository.ExperienceStorer
	ProjectRepo     repository.ProjectStorer
	CertificateRepo repository.CertificateStorer
	AchievementRepo repository.AchievementStorer
}

// Service interface provides methods to interact with user profiles.
type Service interface {
	CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) (int, error)
	ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error)
	GetProfile(ctx context.Context, profileID string) (value dto.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, profileID string, profileDetail dto.UpdateProfileRequest) (ID int, err error)

	UserLoginServive
	EducationService
	ProjectService
	ExperienceService
	CertificateService
	AchievementService
}

type RepoDeps struct {
	UserLoginDeps   repository.UserStorer
	ProfileDeps     repository.ProfileStorer
	EducationDeps   repository.EducationStorer
	ExperienceDeps  repository.ExperienceStorer
	ProjectDeps     repository.ProjectStorer
	CertificateDeps repository.CertificateStorer
	AchievementDeps repository.AchievementStorer
}

// NewServices creates a new instance of the Service.
func NewServices(rp RepoDeps) Service {
	return &service{
		UserLoginRepo:   rp.UserLoginDeps,
		ProfileRepo:     rp.ProfileDeps,
		EducationRepo:   rp.EducationDeps,
		ExperienceRepo:  rp.ExperienceDeps,
		ProjectRepo:     rp.ProjectDeps,
		CertificateRepo: rp.CertificateDeps,
		AchievementRepo: rp.AchievementDeps,
	}
}

var today = helpers.GetTodaysDate()

// CreateProfile : Service layer function creates a new user profile using the provided details.
func (profileSvc *service) CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) (int, error) {
	profileID, err := profileSvc.ProfileRepo.CreateProfile(ctx, profileDetail)
	if err != nil {
		zap.S().Error("Unable to create profile : ", err, " for profile id : ", profileID)
		return 0, err
	}
	return profileID, nil
}

// ListProfiles in the service layer retrieves a list of user profiles.
func (profileSvc *service) ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error) {
	values, err = profileSvc.ProfileRepo.ListProfiles(ctx)
	if err != nil {
		zap.S().Error("Unable to list profile : ", err)
		return []dto.ListProfiles{}, err
	}
	return values, nil
}

// GetProfile in the service layer retrieves a list of user profiles.
func (profileSvc *service) GetProfile(ctx context.Context, profileID string) (value dto.ResponseProfile, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err != nil {
		zap.S().Error("Unable to get profile params : ", err, " for profile id : ", profileID)
		return dto.ResponseProfile{}, err
	}

	value, err = profileSvc.ProfileRepo.GetProfile(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get profile : ", err, " for profile id : ", profileID)
		return dto.ResponseProfile{}, err
	}
	return value, nil
}

// UpdateProfile in the service layer updates user profile.
func (profileSvc *service) UpdateProfile(ctx context.Context, profileID string, profileDetail dto.UpdateProfileRequest) (ID int, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err != nil {
		zap.S().Error("error to get profile params : ", err, " for profile id : ", profileID)
		return 0, err
	}

	ID, err = profileSvc.ProfileRepo.UpdateProfile(ctx, id, profileDetail)
	if err != nil {
		zap.S().Error("Unable to update profile : ", err, " for profile id : ", profileID)
		return 0, err
	}
	return ID, nil
}

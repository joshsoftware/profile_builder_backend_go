package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

// service implements the Service interface.
type service struct {
	ProfileRepo    repository.ProfileStorer	
	EducationRepo    repository.EducationStorer	
	ExperienceRepo repository.ExperienceStorer	
	ProjectRepo repository.ProjectStorer
	CertificateRepo repository.CertificateStorer
	AchievementRepo repository.AchievementStorer
}

// Service interface provides methods to interact with user profiles.
type Service interface {
	CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) (int, error)
	ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error)
	GetProfile(ctx context.Context, profileID string) (value dto.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, profileID string, profileDetail dto.UpdateProfileRequest) (ID int, err error)

	EducationService
	ProjectService
	ExperienceService
	CertificateService
	AchievementService
}

// NewServices creates a new instance of the Service.
func NewServices(db *pgx.Conn) Service {
	profileRepo := repository.NewProfileRepo(db)
	educationRepo := repository.NewEducationRepo(db)
	experienceRepo := repository.NewExperienceRepo(db)
	projectRepo := repository.NewProjectRepo(db)
	certificateRepo := repository.NewCertificateRepo(db)
	achievementRepo := repository.NewAchievementRepo(db)
	return &service{
		ProfileRepo:        profileRepo,
		EducationRepo: educationRepo,
		ExperienceRepo: experienceRepo,
		ProjectRepo: projectRepo,
		CertificateRepo: certificateRepo,
		AchievementRepo: achievementRepo,
	}
}

var today = helpers.GetTodaysDate()

// CreateProfile : Service layer function creates a new user profile using the provided details.
func (profileSvc *service) CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) (int, error) {
	profileID, err := profileSvc.ProfileRepo.CreateProfile(ctx, profileDetail)
	if err != nil {
		return 0, err
	}
	return profileID, nil
}

// ListProfiles in the service layer retrieves a list of user profiles.
func (profileSvc *service) ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error) {
	values, err = profileSvc.ProfileRepo.ListProfiles(ctx)
	if err != nil {
		return []dto.ListProfiles{}, err
	}
	return values, nil
}

// GetProfile in the service layer retrieves a list of user profiles.
func (profileSvc *service) GetProfile(ctx context.Context, profileID string) (value dto.ResponseProfile, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err!= nil {
        return dto.ResponseProfile{}, err
    }

	value, err = profileSvc.ProfileRepo.GetProfile(ctx, id)
	if err != nil {
		return dto.ResponseProfile{}, err
	}
	return value, nil
}

// UpdateProfile in the service layer updates user profile.
func (profileSvc *service) UpdateProfile(ctx context.Context, profileID string, profileDetail dto.UpdateProfileRequest) (ID int, err error) {
	id, err := helpers.ConvertStringToInt(profileID)
	if err!= nil {
        return 0, err
    }

	ID, err = profileSvc.ProfileRepo.UpdateProfile(ctx, id, profileDetail)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

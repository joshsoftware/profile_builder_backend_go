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
	EducationService
	ProjectService
	ExperienceService
	CertificateService
	AchievementService
}

// Service interface provides methods to interact with user profiles.
type Service interface {
	CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) (int, error)
	CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) (profileID int, err error)
	CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (profileID int, err error)
	CreateExperience(ctx context.Context, projDetail dto.CreateExperienceRequest) (profileID int, err error)
	CreateCertificate(ctx context.Context, expDetail dto.CreateCertificateRequest) (profileID int, err error)
	CreateAchievement(ctx context.Context, cDetail dto.CreateAchievementRequest) (profileID int, err error)

	ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error)
}

// NewServices creates a new instance of the Service.
func NewServices(db *pgx.Conn) Service {
	profileRepo := repository.NewProfileRepo(db)
	educationService := NewEducationService(db)
	projectSvc := NewProjectService(db)
	experienceSvc := NewExperienceService(db)
	certificateSvc := NewCertificateService(db)
	achievementSvc := NewAchivementService(db)

	return &service{
		ProfileRepo:        profileRepo,
		EducationService:      educationService,
		ProjectService:        projectSvc,
		ExperienceService:     experienceSvc,
		CertificateService:    certificateSvc,
		AchievementService: achievementSvc,
	}
}

// func NewServices(
// 	profileRepo repository.ProfileStorer,
// 	educationService EducationService,
// 	projectService ProjectService,
// 	experienceService ExperienceService,
// 	certificateService CertificateService,
// 	achievementService AchievementService,
// ) Service {
// 	return &service{
// 		ProfileRepo:        profileRepo,
// 		EducationService:   educationService,
// 		ProjectService:     projectService,
// 		ExperienceService:  experienceService,
// 		CertificateService: certificateService,
// 		AchievementService: achievementService,
// 	}
// }

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

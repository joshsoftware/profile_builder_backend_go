package service

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
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
	CreateProfile(ctx context.Context, profileDetail specs.CreateProfileRequest) (int, error)
	ListProfiles(ctx context.Context) (values []specs.ResponseListProfiles, err error)
	ListSkills(ctx context.Context) (values specs.ListSkills, err error)
	GetProfile(ctx context.Context, id int) (value specs.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, id int, profileDetail specs.UpdateProfileRequest) (ID int, err error)

	UserLoginServive
	EducationService
	ProjectService
	ExperienceService
	CertificateService
	AchievementService
}

// RepoDeps is used to intialize repo dependencies
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
func (profileSvc *service) CreateProfile(ctx context.Context, profileDetail specs.CreateProfileRequest) (int, error) {
	var profileRepo repository.ProfileRepo
	profileRepo.Name = profileDetail.Profile.Name
	profileRepo.Email = profileDetail.Profile.Email
	profileRepo.Gender = profileDetail.Profile.Gender
	profileRepo.Mobile = profileDetail.Profile.Mobile
	profileRepo.Designation = profileDetail.Profile.Designation
	profileRepo.Description = profileDetail.Profile.Description
	profileRepo.Title = profileDetail.Profile.Title
	profileRepo.YearsOfExperience = profileDetail.Profile.YearsOfExperience
	profileRepo.PrimarySkills = profileDetail.Profile.PrimarySkills
	profileRepo.SecondarySkills = profileDetail.Profile.SecondarySkills
	profileRepo.GithubLink = profileDetail.Profile.GithubLink
	profileRepo.LinkedinLink = profileDetail.Profile.LinkedinLink
	profileRepo.CareerObjectives = profileDetail.Profile.CareerObjectives
	profileRepo.CreatedAt = today
	profileRepo.UpdatedAt = today
	//TODO by context
	profileRepo.CreatedByID = 1
	profileRepo.UpdatedByID = 1

	profileID, err := profileSvc.ProfileRepo.CreateProfile(ctx, profileRepo)
	if err != nil {
		zap.S().Error("Unable to create profile : ", err, " for profile id : ", profileID)
		return 0, err
	}
	zap.S().Info("profile created with profile id : ", profileID)

	return profileID, nil
}

// ListProfiles in the service layer retrieves a list of user profiles.
func (profileSvc *service) ListProfiles(ctx context.Context) (values []specs.ResponseListProfiles, err error) {
	profiles, err := profileSvc.ProfileRepo.ListProfiles(ctx)
	if err != nil {
		zap.S().Error("Unable to list profile : ", err)
		return []specs.ResponseListProfiles{}, err
	}

	for _, profile := range profiles {
		isCurrentEmployee := "No"
		if profile.IsCurrentEmployee == 1 {
			isCurrentEmployee = "YES"
		}

		values = append(values, specs.ResponseListProfiles{
			ID:                profile.ID,
			Name:              profile.Name,
			Email:             profile.Email,
			YearsOfExperience: profile.YearsOfExperience,
			PrimarySkills:     profile.PrimarySkills,
			IsCurrentEmployee: isCurrentEmployee,
		})
	}

	return values, nil
}

// ListSkills in the service layer retrieves a list of skills.
func (profileSvc *service) ListSkills(ctx context.Context) (values specs.ListSkills, err error) {
	values, err = profileSvc.ProfileRepo.ListSkills(ctx)
	if err != nil {
		zap.S().Error("Unable to list skills : ", err)
		return specs.ListSkills{}, err
	}
	return values, nil
}

// GetProfile in the service layer retrieves a list of user profiles.
func (profileSvc *service) GetProfile(ctx context.Context, id int) (value specs.ResponseProfile, err error) {

	value, err = profileSvc.ProfileRepo.GetProfile(ctx, id)
	if err != nil {
		zap.S().Error("Unable to get profile : ", err, " for profile id : ", id)
		return specs.ResponseProfile{}, err
	}
	return value, nil
}

// UpdateProfile in the service layer updates user profile.
func (profileSvc *service) UpdateProfile(ctx context.Context, id int, profileDetail specs.UpdateProfileRequest) (ID int, err error) {

	var profileRepo repository.UpdateProfileRepo
	profileRepo.Name = profileDetail.Profile.Name
	profileRepo.Email = profileDetail.Profile.Email
	profileRepo.Gender = profileDetail.Profile.Gender
	profileRepo.Mobile = profileDetail.Profile.Mobile
	profileRepo.Designation = profileDetail.Profile.Designation
	profileRepo.Description = profileDetail.Profile.Description
	profileRepo.Title = profileDetail.Profile.Title
	profileRepo.YearsOfExperience = profileDetail.Profile.YearsOfExperience
	profileRepo.PrimarySkills = profileDetail.Profile.PrimarySkills
	profileRepo.SecondarySkills = profileDetail.Profile.SecondarySkills
	profileRepo.GithubLink = profileDetail.Profile.GithubLink
	profileRepo.LinkedinLink = profileDetail.Profile.LinkedinLink
	profileRepo.UpdatedAt = today
	//TODO by context
	profileRepo.UpdatedByID = 1

	ID, err = profileSvc.ProfileRepo.UpdateProfile(ctx, id, profileRepo)
	if err != nil {
		zap.S().Error("Unable to update profile : ", err, " for profile id : ", id)
		return 0, err
	}
	zap.S().Info("profile update with profile id : ", ID)

	return ID, nil
}

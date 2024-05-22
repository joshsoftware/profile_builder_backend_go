package profile

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

// service implements the Service interface.
type service struct {
	Repo repository.ProfileStorer
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
func NewServices(repo repository.ProfileStorer) Service {
	return &service{
		Repo: repo,
	}
}

var today = helpers.GetTodaysDate()

// CreateProfile : Service layer function creates a new user profile using the provided details.
func (profileSvc *service) CreateProfile(ctx context.Context, profileDetail dto.CreateProfileRequest) (int, error) {

	profileID, err := profileSvc.Repo.CreateProfile(ctx, profileDetail)
	if err != nil {
		return 0, err
	}

	return profileID, nil
}

// CreateEducation : Service layer function adds education details to a user profile.
func (profileSvc *service) CreateEducation(ctx context.Context, eduDetail dto.CreateEducationRequest) (profileID int, err error) {
	if len(eduDetail.Educations) == 0 {
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.EducationDao
	for i := 0; i < len(eduDetail.Educations); i++ {
		var val repository.EducationDao
		val.ProfileID = eduDetail.ProfileID
		val.Degree = eduDetail.Educations[i].Degree
		val.UniversityName = eduDetail.Educations[i].UniversityName
		val.Place = eduDetail.Educations[i].Place
		val.PercentageOrCgpa = eduDetail.Educations[i].PercentageOrCgpa
		val.PassingYear = eduDetail.Educations[i].PassingYear
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.Repo.CreateEducation(ctx, value)
	if err != nil {
		return 0, err
	}

	return eduDetail.ProfileID, nil
}

// CreateProject : Service layer function adds project details to a user profile.
func (profileSvc *service) CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (profileID int, err error) {
	if len(projDetail.Projects) == 0 {
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.ProjectDao
	for i := 0; i < len(projDetail.Projects); i++ {
		var val repository.ProjectDao

		val.ProfileID = projDetail.ProfileID
		val.Name = projDetail.Projects[i].Name
		val.Description = projDetail.Projects[i].Description
		val.Role = projDetail.Projects[i].Role
		val.Responsibilities = projDetail.Projects[i].Responsibilities
		val.Technologies = projDetail.Projects[i].Technologies
		val.TechWorkedOn = projDetail.Projects[i].TechWorkedOn
		val.Duration = projDetail.Projects[i].Duration
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.Repo.CreateProject(ctx, value)
	if err != nil {
		return 0, err
	}

	return projDetail.ProfileID, nil
}

// CreateProject : Service layer function adds experiences details to a user profile.
func (profileSvc *service) CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest) (profileID int, err error) {
	if len(expDetail.Experiences) == 0 {
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.ExperienceDao
	for i := 0; i < len(expDetail.Experiences); i++ {
		var val repository.ExperienceDao

		val.ProfileID = expDetail.ProfileID
		val.Designation = expDetail.Experiences[i].Designation
		val.CompanyName = expDetail.Experiences[i].CompanyName
		val.FromDate = expDetail.Experiences[i].FromDate
		val.ToDate = expDetail.Experiences[i].ToDate
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.Repo.CreateExperience(ctx, value)
	if err != nil {
		return 0, err
	}

	return expDetail.ProfileID, nil
}

// CreateCerticate : Service layer function adds certicates details to a user profile.
func (profileSvc *service) CreateCertificate(ctx context.Context, cDetail dto.CreateCertificateRequest) (profileID int, err error) {
	if len(cDetail.Certificates) == 0 {
		return 0, errors.ErrEmptyPayload
	}

	var value []repository.CertificateDao
	for i := 0; i < len(cDetail.Certificates); i++ {
		var val repository.CertificateDao

		val.ProfileID = cDetail.ProfileID
		val.Name = cDetail.Certificates[i].Name
		val.OrganizationName = cDetail.Certificates[i].OrganizationName
		val.Description = cDetail.Certificates[i].Description
		val.IssuedDate = cDetail.Certificates[i].IssuedDate
		val.FromDate = cDetail.Certificates[i].FromDate
		val.ToDate = cDetail.Certificates[i].ToDate
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.Repo.CreateCertificate(ctx, value)
	if err != nil {
		return 0, err
	}

	return cDetail.ProfileID, nil
}

// CreateAchievement : Service layer function adds certicates details to a user profile.
func (profileSvc *service) CreateAchievement(ctx context.Context, cDetail dto.CreateAchievementRequest) (profileID int, err error) {
	if len(cDetail.Achievements) == 0 {
		return 0, errors.ErrEmptyPayload
	}
	var value []repository.AchievementDao
	for i := 0; i < len(cDetail.Achievements); i++ {
		var val repository.AchievementDao

		val.ProfileID = cDetail.ProfileID
		val.Name = cDetail.Achievements[i].Name
		val.Description = cDetail.Achievements[i].Description
		val.CreatedAt = today
		val.UpdatedAt = today
		val.CreatedByID = 1
		val.UpdatedByID = 1

		value = append(value, val)
	}

	err = profileSvc.Repo.CreateAchievement(ctx, value)
	if err != nil {
		return 0, err
	}

	return cDetail.ProfileID, nil
}

// ListProfiles in the service layer retrieves a list of user profiles.
func (profileSvc *service) ListProfiles(ctx context.Context) (values []dto.ListProfiles, err error) {

	values, err = profileSvc.Repo.ListProfiles(ctx)
	if err != nil {
		return []dto.ListProfiles{}, err
	}

	return values, nil
}

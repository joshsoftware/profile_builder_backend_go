// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// BackupAllProfiles provides a mock function with given fields:
func (_m *Service) BackupAllProfiles() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for BackupAllProfiles")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAchievement provides a mock function with given fields: ctx, cDetail, profileID, userID
func (_m *Service) CreateAchievement(ctx context.Context, cDetail specs.CreateAchievementRequest, profileID int, userID int) (int, error) {
	ret := _m.Called(ctx, cDetail, profileID, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateAchievement")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateAchievementRequest, int, int) (int, error)); ok {
		return rf(ctx, cDetail, profileID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateAchievementRequest, int, int) int); ok {
		r0 = rf(ctx, cDetail, profileID, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateAchievementRequest, int, int) error); ok {
		r1 = rf(ctx, cDetail, profileID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCertificate provides a mock function with given fields: ctx, cDetail, profileID, userID
func (_m *Service) CreateCertificate(ctx context.Context, cDetail specs.CreateCertificateRequest, profileID int, userID int) (int, error) {
	ret := _m.Called(ctx, cDetail, profileID, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateCertificate")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateCertificateRequest, int, int) (int, error)); ok {
		return rf(ctx, cDetail, profileID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateCertificateRequest, int, int) int); ok {
		r0 = rf(ctx, cDetail, profileID, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateCertificateRequest, int, int) error); ok {
		r1 = rf(ctx, cDetail, profileID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateEducation provides a mock function with given fields: ctx, eduDetail, profileID, userID
func (_m *Service) CreateEducation(ctx context.Context, eduDetail specs.CreateEducationRequest, profileID int, userID int) (int, error) {
	ret := _m.Called(ctx, eduDetail, profileID, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateEducation")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateEducationRequest, int, int) (int, error)); ok {
		return rf(ctx, eduDetail, profileID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateEducationRequest, int, int) int); ok {
		r0 = rf(ctx, eduDetail, profileID, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateEducationRequest, int, int) error); ok {
		r1 = rf(ctx, eduDetail, profileID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateExperience provides a mock function with given fields: ctx, expDetail, profileID, userID
func (_m *Service) CreateExperience(ctx context.Context, expDetail specs.CreateExperienceRequest, profileID int, userID int) (int, error) {
	ret := _m.Called(ctx, expDetail, profileID, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateExperienceRequest, int, int) (int, error)); ok {
		return rf(ctx, expDetail, profileID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateExperienceRequest, int, int) int); ok {
		r0 = rf(ctx, expDetail, profileID, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateExperienceRequest, int, int) error); ok {
		r1 = rf(ctx, expDetail, profileID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProfile provides a mock function with given fields: ctx, profileDetail, userID
func (_m *Service) CreateProfile(ctx context.Context, profileDetail specs.CreateProfileRequest, userID int) (int, error) {
	ret := _m.Called(ctx, profileDetail, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateProfile")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateProfileRequest, int) (int, error)); ok {
		return rf(ctx, profileDetail, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateProfileRequest, int) int); ok {
		r0 = rf(ctx, profileDetail, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateProfileRequest, int) error); ok {
		r1 = rf(ctx, profileDetail, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProject provides a mock function with given fields: ctx, projDetail, profileID, userID
func (_m *Service) CreateProject(ctx context.Context, projDetail specs.CreateProjectRequest, profileID int, userID int) (int, error) {
	ret := _m.Called(ctx, projDetail, profileID, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateProject")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateProjectRequest, int, int) (int, error)); ok {
		return rf(ctx, projDetail, profileID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateProjectRequest, int, int) int); ok {
		r0 = rf(ctx, projDetail, profileID, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateProjectRequest, int, int) error); ok {
		r1 = rf(ctx, projDetail, profileID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAchievement provides a mock function with given fields: ctx, profileID, achievementID
func (_m *Service) DeleteAchievement(ctx context.Context, profileID int, achievementID int) error {
	ret := _m.Called(ctx, profileID, achievementID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAchievement")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, profileID, achievementID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCertificate provides a mock function with given fields: ctx, profileID, certitificateID
func (_m *Service) DeleteCertificate(ctx context.Context, profileID int, certitificateID int) error {
	ret := _m.Called(ctx, profileID, certitificateID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCertificate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, profileID, certitificateID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEducation provides a mock function with given fields: ctx, profileID, educationID
func (_m *Service) DeleteEducation(ctx context.Context, profileID int, educationID int) error {
	ret := _m.Called(ctx, profileID, educationID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEducation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, profileID, educationID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteExperience provides a mock function with given fields: ctx, profileID, experienceID
func (_m *Service) DeleteExperience(ctx context.Context, profileID int, experienceID int) error {
	ret := _m.Called(ctx, profileID, experienceID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteExperience")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, profileID, experienceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteProfile provides a mock function with given fields: ctx, profileID
func (_m *Service) DeleteProfile(ctx context.Context, profileID int) error {
	ret := _m.Called(ctx, profileID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProfile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, profileID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteProject provides a mock function with given fields: ctx, profileID, projectID
func (_m *Service) DeleteProject(ctx context.Context, profileID int, projectID int) error {
	ret := _m.Called(ctx, profileID, projectID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, profileID, projectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateLoginToken provides a mock function with given fields: ctx, email
func (_m *Service) GenerateLoginToken(ctx context.Context, email string) (string, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GenerateLoginToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfile provides a mock function with given fields: ctx, id
func (_m *Service) GetProfile(ctx context.Context, id int) (specs.ResponseProfile, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetProfile")
	}

	var r0 specs.ResponseProfile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (specs.ResponseProfile, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) specs.ResponseProfile); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(specs.ResponseProfile)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAchievements provides a mock function with given fields: ctx, profileID, filter
func (_m *Service) ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter) ([]specs.AchievementResponse, error) {
	ret := _m.Called(ctx, profileID, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListAchievements")
	}

	var r0 []specs.AchievementResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListAchievementFilter) ([]specs.AchievementResponse, error)); ok {
		return rf(ctx, profileID, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListAchievementFilter) []specs.AchievementResponse); ok {
		r0 = rf(ctx, profileID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.AchievementResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListAchievementFilter) error); ok {
		r1 = rf(ctx, profileID, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListCertificates provides a mock function with given fields: ctx, profileID, fitler
func (_m *Service) ListCertificates(ctx context.Context, profileID int, fitler specs.ListCertificateFilter) ([]specs.CertificateResponse, error) {
	ret := _m.Called(ctx, profileID, fitler)

	if len(ret) == 0 {
		panic("no return value specified for ListCertificates")
	}

	var r0 []specs.CertificateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListCertificateFilter) ([]specs.CertificateResponse, error)); ok {
		return rf(ctx, profileID, fitler)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListCertificateFilter) []specs.CertificateResponse); ok {
		r0 = rf(ctx, profileID, fitler)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.CertificateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListCertificateFilter) error); ok {
		r1 = rf(ctx, profileID, fitler)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListEducations provides a mock function with given fields: ctx, id, filter
func (_m *Service) ListEducations(ctx context.Context, id int, filter specs.ListEducationsFilter) ([]specs.EducationResponse, error) {
	ret := _m.Called(ctx, id, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListEducations")
	}

	var r0 []specs.EducationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListEducationsFilter) ([]specs.EducationResponse, error)); ok {
		return rf(ctx, id, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListEducationsFilter) []specs.EducationResponse); ok {
		r0 = rf(ctx, id, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.EducationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListEducationsFilter) error); ok {
		r1 = rf(ctx, id, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListExperiences provides a mock function with given fields: ctx, id, filter
func (_m *Service) ListExperiences(ctx context.Context, id int, filter specs.ListExperiencesFilter) ([]specs.ExperienceResponse, error) {
	ret := _m.Called(ctx, id, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListExperiences")
	}

	var r0 []specs.ExperienceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListExperiencesFilter) ([]specs.ExperienceResponse, error)); ok {
		return rf(ctx, id, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListExperiencesFilter) []specs.ExperienceResponse); ok {
		r0 = rf(ctx, id, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ExperienceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListExperiencesFilter) error); ok {
		r1 = rf(ctx, id, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProfiles provides a mock function with given fields: ctx
func (_m *Service) ListProfiles(ctx context.Context) ([]specs.ResponseListProfiles, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListProfiles")
	}

	var r0 []specs.ResponseListProfiles
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]specs.ResponseListProfiles, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []specs.ResponseListProfiles); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ResponseListProfiles)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProjects provides a mock function with given fields: ctx, profileID, filter
func (_m *Service) ListProjects(ctx context.Context, profileID int, filter specs.ListProjectsFilter) ([]specs.ProjectResponse, error) {
	ret := _m.Called(ctx, profileID, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListProjects")
	}

	var r0 []specs.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListProjectsFilter) ([]specs.ProjectResponse, error)); ok {
		return rf(ctx, profileID, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListProjectsFilter) []specs.ProjectResponse); ok {
		r0 = rf(ctx, profileID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListProjectsFilter) error); ok {
		r1 = rf(ctx, profileID, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSkills provides a mock function with given fields: ctx
func (_m *Service) ListSkills(ctx context.Context) (specs.ListSkills, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListSkills")
	}

	var r0 specs.ListSkills
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (specs.ListSkills, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) specs.ListSkills); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(specs.ListSkills)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAchievement provides a mock function with given fields: ctx, profileID, achID, userID, req
func (_m *Service) UpdateAchievement(ctx context.Context, profileID int, achID int, userID int, req specs.UpdateAchievementRequest) (int, error) {
	ret := _m.Called(ctx, profileID, achID, userID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAchievement")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateAchievementRequest) (int, error)); ok {
		return rf(ctx, profileID, achID, userID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateAchievementRequest) int); ok {
		r0 = rf(ctx, profileID, achID, userID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int, specs.UpdateAchievementRequest) error); ok {
		r1 = rf(ctx, profileID, achID, userID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCertificate provides a mock function with given fields: ctx, profileID, certID, userID, req
func (_m *Service) UpdateCertificate(ctx context.Context, profileID int, certID int, userID int, req specs.UpdateCertificateRequest) (int, error) {
	ret := _m.Called(ctx, profileID, certID, userID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCertificate")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateCertificateRequest) (int, error)); ok {
		return rf(ctx, profileID, certID, userID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateCertificateRequest) int); ok {
		r0 = rf(ctx, profileID, certID, userID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int, specs.UpdateCertificateRequest) error); ok {
		r1 = rf(ctx, profileID, certID, userID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEducation provides a mock function with given fields: ctx, profileID, eduID, userID, req
func (_m *Service) UpdateEducation(ctx context.Context, profileID int, eduID int, userID int, req specs.UpdateEducationRequest) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, userID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateEducation")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateEducationRequest) (int, error)); ok {
		return rf(ctx, profileID, eduID, userID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateEducationRequest) int); ok {
		r0 = rf(ctx, profileID, eduID, userID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int, specs.UpdateEducationRequest) error); ok {
		r1 = rf(ctx, profileID, eduID, userID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateExperience provides a mock function with given fields: ctx, profileID, expID, userID, req
func (_m *Service) UpdateExperience(ctx context.Context, profileID int, expID int, userID int, req specs.UpdateExperienceRequest) (int, error) {
	ret := _m.Called(ctx, profileID, expID, userID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateExperienceRequest) (int, error)); ok {
		return rf(ctx, profileID, expID, userID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateExperienceRequest) int); ok {
		r0 = rf(ctx, profileID, expID, userID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int, specs.UpdateExperienceRequest) error); ok {
		r1 = rf(ctx, profileID, expID, userID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfile provides a mock function with given fields: ctx, profileID, userID, profileDetail
func (_m *Service) UpdateProfile(ctx context.Context, profileID int, userID int, profileDetail specs.UpdateProfileRequest) (int, error) {
	ret := _m.Called(ctx, profileID, userID, profileDetail)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProfile")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, specs.UpdateProfileRequest) (int, error)); ok {
		return rf(ctx, profileID, userID, profileDetail)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, specs.UpdateProfileRequest) int); ok {
		r0 = rf(ctx, profileID, userID, profileDetail)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, specs.UpdateProfileRequest) error); ok {
		r1 = rf(ctx, profileID, userID, profileDetail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfileStatus provides a mock function with given fields: ctx, profileID, req
func (_m *Service) UpdateProfileStatus(ctx context.Context, profileID int, req specs.UpdateProfileStatus) error {
	ret := _m.Called(ctx, profileID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProfileStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.UpdateProfileStatus) error); ok {
		r0 = rf(ctx, profileID, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateProject provides a mock function with given fields: ctx, profileID, projID, userID, req
func (_m *Service) UpdateProject(ctx context.Context, profileID int, projID int, userID int, req specs.UpdateProjectRequest) (int, error) {
	ret := _m.Called(ctx, profileID, projID, userID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProject")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateProjectRequest) (int, error)); ok {
		return rf(ctx, profileID, projID, userID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int, specs.UpdateProjectRequest) int); ok {
		r0 = rf(ctx, profileID, projID, userID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int, specs.UpdateProjectRequest) error); ok {
		r1 = rf(ctx, profileID, projID, userID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSequence provides a mock function with given fields: ctx, userID, seqDetail
func (_m *Service) UpdateSequence(ctx context.Context, userID int, seqDetail specs.UpdateSequenceRequest) (int, error) {
	ret := _m.Called(ctx, userID, seqDetail)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSequence")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.UpdateSequenceRequest) (int, error)); ok {
		return rf(ctx, userID, seqDetail)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.UpdateSequenceRequest) int); ok {
		r0 = rf(ctx, userID, seqDetail)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.UpdateSequenceRequest) error); ok {
		r1 = rf(ctx, userID, seqDetail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

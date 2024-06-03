// Code generated by mockery v2.40.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	mock "github.com/stretchr/testify/mock"
)

// ExperienceService is an autogenerated mock type for the ExperienceService type
type ExperienceService struct {
	mock.Mock
}

// CreateExperience provides a mock function with given fields: ctx, expDetail, ID
func (_m *ExperienceService) CreateExperience(ctx context.Context, expDetail dto.CreateExperienceRequest, ID string) (int, error) {
	ret := _m.Called(ctx, expDetail, ID)

	if len(ret) == 0 {
		panic("no return value specified for CreateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateExperienceRequest, string) (int, error)); ok {
		return rf(ctx, expDetail, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateExperienceRequest, string) int); ok {
		r0 = rf(ctx, expDetail, ID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.CreateExperienceRequest, string) error); ok {
		r1 = rf(ctx, expDetail, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExperience provides a mock function with given fields: ctx, profileID
func (_m *ExperienceService) GetExperience(ctx context.Context, profileID string) ([]dto.ExperienceResponse, error) {
	ret := _m.Called(ctx, profileID)

	if len(ret) == 0 {
		panic("no return value specified for GetExperience")
	}

	var r0 []dto.ExperienceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]dto.ExperienceResponse, error)); ok {
		return rf(ctx, profileID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []dto.ExperienceResponse); ok {
		r0 = rf(ctx, profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.ExperienceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, profileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateExperience provides a mock function with given fields: ctx, profileID, eduID, req
func (_m *ExperienceService) UpdateExperience(ctx context.Context, profileID string, eduID string, req dto.UpdateExperienceRequest) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, dto.UpdateExperienceRequest) (int, error)); ok {
		return rf(ctx, profileID, eduID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, dto.UpdateExperienceRequest) int); ok {
		r0 = rf(ctx, profileID, eduID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, dto.UpdateExperienceRequest) error); ok {
		r1 = rf(ctx, profileID, eduID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExperienceService creates a new instance of ExperienceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExperienceService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExperienceService {
	mock := &ExperienceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

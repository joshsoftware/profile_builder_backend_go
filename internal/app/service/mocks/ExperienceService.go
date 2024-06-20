// Code generated by mockery v2.40.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// ExperienceService is an autogenerated mock type for the ExperienceService type
type ExperienceService struct {
	mock.Mock
}

// CreateExperience provides a mock function with given fields: ctx, expDetail, id
func (_m *ExperienceService) CreateExperience(ctx context.Context, expDetail specs.CreateExperienceRequest, id int) (int, error) {
	ret := _m.Called(ctx, expDetail, id)

	if len(ret) == 0 {
		panic("no return value specified for CreateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateExperienceRequest, int) (int, error)); ok {
		return rf(ctx, expDetail, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, specs.CreateExperienceRequest, int) int); ok {
		r0 = rf(ctx, expDetail, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, specs.CreateExperienceRequest, int) error); ok {
		r1 = rf(ctx, expDetail, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExperience provides a mock function with given fields: ctx, id
func (_m *ExperienceService) GetExperience(ctx context.Context, id int) ([]specs.ExperienceResponse, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetExperience")
	}

	var r0 []specs.ExperienceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]specs.ExperienceResponse, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []specs.ExperienceResponse); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ExperienceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateExperience provides a mock function with given fields: ctx, profileID, eduID, req
func (_m *ExperienceService) UpdateExperience(ctx context.Context, profileID string, eduID string, req specs.UpdateExperienceRequest) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, specs.UpdateExperienceRequest) (int, error)); ok {
		return rf(ctx, profileID, eduID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, specs.UpdateExperienceRequest) int); ok {
		r0 = rf(ctx, profileID, eduID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, specs.UpdateExperienceRequest) error); ok {
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

// Code generated by mockery v2.40.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	mock "github.com/stretchr/testify/mock"
)

// ProjectService is an autogenerated mock type for the ProjectService type
type ProjectService struct {
	mock.Mock
}

// CreateProject provides a mock function with given fields: ctx, projDetail
func (_m *ProjectService) CreateProject(ctx context.Context, projDetail dto.CreateProjectRequest) (int, error) {
	ret := _m.Called(ctx, projDetail)

	if len(ret) == 0 {
		panic("no return value specified for CreateProject")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateProjectRequest) (int, error)); ok {
		return rf(ctx, projDetail)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateProjectRequest) int); ok {
		r0 = rf(ctx, projDetail)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.CreateProjectRequest) error); ok {
		r1 = rf(ctx, projDetail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProject provides a mock function with given fields: ctx, profileID
func (_m *ProjectService) GetProject(ctx context.Context, profileID string) ([]dto.ProjectResponse, error) {
	ret := _m.Called(ctx, profileID)

	if len(ret) == 0 {
		panic("no return value specified for GetProject")
	}

	var r0 []dto.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]dto.ProjectResponse, error)); ok {
		return rf(ctx, profileID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []dto.ProjectResponse); ok {
		r0 = rf(ctx, profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, profileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectService creates a new instance of ProjectService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectService {
	mock := &ProjectService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
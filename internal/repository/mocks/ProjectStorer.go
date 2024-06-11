// Code generated by mockery v2.40.0. DO NOT EDIT.

package mocks

import (
	context "context"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// ProjectStorer is an autogenerated mock type for the ProjectStorer type
type ProjectStorer struct {
	mock.Mock
}

// CreateProject provides a mock function with given fields: ctx, values
func (_m *ProjectStorer) CreateProject(ctx context.Context, values []repository.ProjectRepo) error {
	ret := _m.Called(ctx, values)

	if len(ret) == 0 {
		panic("no return value specified for CreateProject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []repository.ProjectRepo) error); ok {
		r0 = rf(ctx, values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProjects provides a mock function with given fields: ctx, profileID
func (_m *ProjectStorer) GetProjects(ctx context.Context, profileID int) ([]specs.ProjectResponse, error) {
	ret := _m.Called(ctx, profileID)

	if len(ret) == 0 {
		panic("no return value specified for GetProjects")
	}

	var r0 []specs.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]specs.ProjectResponse, error)); ok {
		return rf(ctx, profileID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []specs.ProjectResponse); ok {
		r0 = rf(ctx, profileID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, profileID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProject provides a mock function with given fields: ctx, profileID, eduID, req
func (_m *ProjectStorer) UpdateProject(ctx context.Context, profileID int, eduID int, req repository.UpdateProjectRepo) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProject")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateProjectRepo) (int, error)); ok {
		return rf(ctx, profileID, eduID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateProjectRepo) int); ok {
		r0 = rf(ctx, profileID, eduID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, repository.UpdateProjectRepo) error); ok {
		r1 = rf(ctx, profileID, eduID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProjectStorer creates a new instance of ProjectStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProjectStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProjectStorer {
	mock := &ProjectStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

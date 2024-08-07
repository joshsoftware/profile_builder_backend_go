// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// ProjectStorer is an autogenerated mock type for the ProjectStorer type
type ProjectStorer struct {
	mock.Mock
}

// CreateProject provides a mock function with given fields: ctx, values, tx
func (_m *ProjectStorer) CreateProject(ctx context.Context, values []repository.ProjectRepo, tx pgx.Tx) error {
	ret := _m.Called(ctx, values, tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateProject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []repository.ProjectRepo, pgx.Tx) error); ok {
		r0 = rf(ctx, values, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteProject provides a mock function with given fields: ctx, profileID, projectID, tx
func (_m *ProjectStorer) DeleteProject(ctx context.Context, profileID int, projectID int, tx pgx.Tx) error {
	ret := _m.Called(ctx, profileID, projectID, tx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, pgx.Tx) error); ok {
		r0 = rf(ctx, profileID, projectID, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListProjects provides a mock function with given fields: ctx, profileID, filter, tx
func (_m *ProjectStorer) ListProjects(ctx context.Context, profileID int, filter specs.ListProjectsFilter, tx pgx.Tx) ([]specs.ProjectResponse, error) {
	ret := _m.Called(ctx, profileID, filter, tx)

	if len(ret) == 0 {
		panic("no return value specified for ListProjects")
	}

	var r0 []specs.ProjectResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListProjectsFilter, pgx.Tx) ([]specs.ProjectResponse, error)); ok {
		return rf(ctx, profileID, filter, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListProjectsFilter, pgx.Tx) []specs.ProjectResponse); ok {
		r0 = rf(ctx, profileID, filter, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ProjectResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListProjectsFilter, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, filter, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProject provides a mock function with given fields: ctx, profileID, eduID, req, tx
func (_m *ProjectStorer) UpdateProject(ctx context.Context, profileID int, eduID int, req repository.UpdateProjectRepo, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, req, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProject")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateProjectRepo, pgx.Tx) (int, error)); ok {
		return rf(ctx, profileID, eduID, req, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateProjectRepo, pgx.Tx) int); ok {
		r0 = rf(ctx, profileID, eduID, req, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, repository.UpdateProjectRepo, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, eduID, req, tx)
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

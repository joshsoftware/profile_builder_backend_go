// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// ExperienceStorer is an autogenerated mock type for the ExperienceStorer type
type ExperienceStorer struct {
	mock.Mock
}

// CreateExperience provides a mock function with given fields: ctx, values, tx
func (_m *ExperienceStorer) CreateExperience(ctx context.Context, values []repository.ExperienceRepo, tx pgx.Tx) error {
	ret := _m.Called(ctx, values, tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateExperience")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []repository.ExperienceRepo, pgx.Tx) error); ok {
		r0 = rf(ctx, values, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteExperience provides a mock function with given fields: ctx, profileID, experienceID, tx
func (_m *ExperienceStorer) DeleteExperience(ctx context.Context, profileID int, experienceID int, tx pgx.Tx) error {
	ret := _m.Called(ctx, profileID, experienceID, tx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteExperience")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, pgx.Tx) error); ok {
		r0 = rf(ctx, profileID, experienceID, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListExperiences provides a mock function with given fields: ctx, profileID, filter, tx
func (_m *ExperienceStorer) ListExperiences(ctx context.Context, profileID int, filter specs.ListExperiencesFilter, tx pgx.Tx) ([]specs.ExperienceResponse, error) {
	ret := _m.Called(ctx, profileID, filter, tx)

	if len(ret) == 0 {
		panic("no return value specified for ListExperiences")
	}

	var r0 []specs.ExperienceResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListExperiencesFilter, pgx.Tx) ([]specs.ExperienceResponse, error)); ok {
		return rf(ctx, profileID, filter, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListExperiencesFilter, pgx.Tx) []specs.ExperienceResponse); ok {
		r0 = rf(ctx, profileID, filter, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ExperienceResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListExperiencesFilter, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, filter, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateExperience provides a mock function with given fields: ctx, profileID, eduID, req, tx
func (_m *ExperienceStorer) UpdateExperience(ctx context.Context, profileID int, eduID int, req repository.UpdateExperienceRepo, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, req, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateExperience")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateExperienceRepo, pgx.Tx) (int, error)); ok {
		return rf(ctx, profileID, eduID, req, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateExperienceRepo, pgx.Tx) int); ok {
		r0 = rf(ctx, profileID, eduID, req, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, repository.UpdateExperienceRepo, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, eduID, req, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExperienceStorer creates a new instance of ExperienceStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExperienceStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExperienceStorer {
	mock := &ExperienceStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

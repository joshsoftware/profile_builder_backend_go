// Code generated by mockery v2.40.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	mock "github.com/stretchr/testify/mock"
)

// ProfileStorer is an autogenerated mock type for the ProfileStorer type
type ProfileStorer struct {
	mock.Mock
}

// CreateProfile provides a mock function with given fields: ctx, pd
func (_m *ProfileStorer) CreateProfile(ctx context.Context, pd dto.CreateProfileRequest) (int, error) {
	ret := _m.Called(ctx, pd)

	if len(ret) == 0 {
		panic("no return value specified for CreateProfile")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateProfileRequest) (int, error)); ok {
		return rf(ctx, pd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.CreateProfileRequest) int); ok {
		r0 = rf(ctx, pd)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.CreateProfileRequest) error); ok {
		r1 = rf(ctx, pd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProfiles provides a mock function with given fields: ctx
func (_m *ProfileStorer) ListProfiles(ctx context.Context) ([]dto.ListProfiles, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListProfiles")
	}

	var r0 []dto.ListProfiles
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]dto.ListProfiles, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []dto.ListProfiles); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.ListProfiles)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProfileStorer creates a new instance of ProfileStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProfileStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProfileStorer {
	mock := &ProfileStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

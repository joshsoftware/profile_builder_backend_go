// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// AchievementService is an autogenerated mock type for the AchievementService type
type AchievementService struct {
	mock.Mock
}

// CreateAchievement provides a mock function with given fields: ctx, cDetail, profileID, userID
func (_m *AchievementService) CreateAchievement(ctx context.Context, cDetail specs.CreateAchievementRequest, profileID int, userID int) (int, error) {
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

// DeleteAchievement provides a mock function with given fields: ctx, req
func (_m *AchievementService) DeleteAchievement(ctx context.Context, req specs.DeleteAchievementRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAchievement")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, specs.DeleteAchievementRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListAchievements provides a mock function with given fields: ctx, profileID, filter
func (_m *AchievementService) ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter) ([]specs.AchievementResponse, error) {
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

// UpdateAchievement provides a mock function with given fields: ctx, profileID, achID, userID, req
func (_m *AchievementService) UpdateAchievement(ctx context.Context, profileID int, achID int, userID int, req specs.UpdateAchievementRequest) (int, error) {
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

// NewAchievementService creates a new instance of AchievementService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAchievementService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AchievementService {
	mock := &AchievementService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

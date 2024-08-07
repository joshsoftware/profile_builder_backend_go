// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// AchievementStorer is an autogenerated mock type for the AchievementStorer type
type AchievementStorer struct {
	mock.Mock
}

// CreateAchievement provides a mock function with given fields: ctx, values, tx
func (_m *AchievementStorer) CreateAchievement(ctx context.Context, values []repository.AchievementRepo, tx pgx.Tx) error {
	ret := _m.Called(ctx, values, tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateAchievement")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []repository.AchievementRepo, pgx.Tx) error); ok {
		r0 = rf(ctx, values, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAchievement provides a mock function with given fields: ctx, profileID, achievementID, tx
func (_m *AchievementStorer) DeleteAchievement(ctx context.Context, profileID int, achievementID int, tx pgx.Tx) error {
	ret := _m.Called(ctx, profileID, achievementID, tx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAchievement")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, pgx.Tx) error); ok {
		r0 = rf(ctx, profileID, achievementID, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListAchievements provides a mock function with given fields: ctx, profileID, filter, tx
func (_m *AchievementStorer) ListAchievements(ctx context.Context, profileID int, filter specs.ListAchievementFilter, tx pgx.Tx) ([]specs.AchievementResponse, error) {
	ret := _m.Called(ctx, profileID, filter, tx)

	if len(ret) == 0 {
		panic("no return value specified for ListAchievements")
	}

	var r0 []specs.AchievementResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListAchievementFilter, pgx.Tx) ([]specs.AchievementResponse, error)); ok {
		return rf(ctx, profileID, filter, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListAchievementFilter, pgx.Tx) []specs.AchievementResponse); ok {
		r0 = rf(ctx, profileID, filter, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.AchievementResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListAchievementFilter, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, filter, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAchievement provides a mock function with given fields: ctx, profileID, achID, req, tx
func (_m *AchievementStorer) UpdateAchievement(ctx context.Context, profileID int, achID int, req repository.UpdateAchievementRepo, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, profileID, achID, req, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAchievement")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateAchievementRepo, pgx.Tx) (int, error)); ok {
		return rf(ctx, profileID, achID, req, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateAchievementRepo, pgx.Tx) int); ok {
		r0 = rf(ctx, profileID, achID, req, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, repository.UpdateAchievementRepo, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, achID, req, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAchievementStorer creates a new instance of AchievementStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAchievementStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *AchievementStorer {
	mock := &AchievementStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

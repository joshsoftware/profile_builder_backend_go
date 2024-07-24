// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// ProfileStorer is an autogenerated mock type for the ProfileStorer type
type ProfileStorer struct {
	mock.Mock
}

// BackupAllProfiles provides a mock function with given fields: backup_dir
func (_m *ProfileStorer) BackupAllProfiles(backup_dir string) {
	_m.Called(backup_dir)
}

// BeginTransaction provides a mock function with given fields: ctx
func (_m *ProfileStorer) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for BeginTransaction")
	}

	var r0 pgx.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (pgx.Tx, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) pgx.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountRecords provides a mock function with given fields: ctx, ProfileID, ComponentName, tx
func (_m *ProfileStorer) CountRecords(ctx context.Context, ProfileID int, ComponentName string, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, ProfileID, ComponentName, tx)

	if len(ret) == 0 {
		panic("no return value specified for CountRecords")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string, pgx.Tx) (int, error)); ok {
		return rf(ctx, ProfileID, ComponentName, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string, pgx.Tx) int); ok {
		r0 = rf(ctx, ProfileID, ComponentName, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string, pgx.Tx) error); ok {
		r1 = rf(ctx, ProfileID, ComponentName, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateProfile provides a mock function with given fields: ctx, pd, tx
func (_m *ProfileStorer) CreateProfile(ctx context.Context, pd repository.ProfileRepo, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, pd, tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateProfile")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.ProfileRepo, pgx.Tx) (int, error)); ok {
		return rf(ctx, pd, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.ProfileRepo, pgx.Tx) int); ok {
		r0 = rf(ctx, pd, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.ProfileRepo, pgx.Tx) error); ok {
		r1 = rf(ctx, pd, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProfile provides a mock function with given fields: ctx, profileID, tx
func (_m *ProfileStorer) DeleteProfile(ctx context.Context, profileID int, tx pgx.Tx) error {
	ret := _m.Called(ctx, profileID, tx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProfile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, pgx.Tx) error); ok {
		r0 = rf(ctx, profileID, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProfile provides a mock function with given fields: ctx, profileID, tx
func (_m *ProfileStorer) GetProfile(ctx context.Context, profileID int, tx pgx.Tx) (specs.ResponseProfile, error) {
	ret := _m.Called(ctx, profileID, tx)

	if len(ret) == 0 {
		panic("no return value specified for GetProfile")
	}

	var r0 specs.ResponseProfile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, pgx.Tx) (specs.ResponseProfile, error)); ok {
		return rf(ctx, profileID, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, pgx.Tx) specs.ResponseProfile); ok {
		r0 = rf(ctx, profileID, tx)
	} else {
		r0 = ret.Get(0).(specs.ResponseProfile)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleTransaction provides a mock function with given fields: ctx, tx, incomingErr
func (_m *ProfileStorer) HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) error {
	ret := _m.Called(ctx, tx, incomingErr)

	if len(ret) == 0 {
		panic("no return value specified for HandleTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx, error) error); ok {
		r0 = rf(ctx, tx, incomingErr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListProfiles provides a mock function with given fields: ctx, tx
func (_m *ProfileStorer) ListProfiles(ctx context.Context, tx pgx.Tx) ([]specs.ListProfiles, error) {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for ListProfiles")
	}

	var r0 []specs.ListProfiles
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) ([]specs.ListProfiles, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) []specs.ListProfiles); ok {
		r0 = rf(ctx, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.ListProfiles)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSkills provides a mock function with given fields: ctx, tx
func (_m *ProfileStorer) ListSkills(ctx context.Context, tx pgx.Tx) (specs.ListSkills, error) {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for ListSkills")
	}

	var r0 specs.ListSkills
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) (specs.ListSkills, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) specs.ListSkills); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Get(0).(specs.ListSkills)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfile provides a mock function with given fields: ctx, profileID, pd, tx
func (_m *ProfileStorer) UpdateProfile(ctx context.Context, profileID int, pd repository.UpdateProfileRepo, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, profileID, pd, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProfile")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, repository.UpdateProfileRepo, pgx.Tx) (int, error)); ok {
		return rf(ctx, profileID, pd, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, repository.UpdateProfileRepo, pgx.Tx) int); ok {
		r0 = rf(ctx, profileID, pd, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, repository.UpdateProfileRepo, pgx.Tx) error); ok {
		r1 = rf(ctx, profileID, pd, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProfileStatus provides a mock function with given fields: ctx, profileID, updateRequest, tx
func (_m *ProfileStorer) UpdateProfileStatus(ctx context.Context, profileID int, updateRequest repository.UpdateProfileStatusRepo, tx pgx.Tx) error {
	ret := _m.Called(ctx, profileID, updateRequest, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProfileStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, repository.UpdateProfileStatusRepo, pgx.Tx) error); ok {
		r0 = rf(ctx, profileID, updateRequest, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSequence provides a mock function with given fields: ctx, us, tx
func (_m *ProfileStorer) UpdateSequence(ctx context.Context, us repository.UpdateSequenceRequest, tx pgx.Tx) (int, error) {
	ret := _m.Called(ctx, us, tx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSequence")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.UpdateSequenceRequest, pgx.Tx) (int, error)); ok {
		return rf(ctx, us, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.UpdateSequenceRequest, pgx.Tx) int); ok {
		r0 = rf(ctx, us, tx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.UpdateSequenceRequest, pgx.Tx) error); ok {
		r1 = rf(ctx, us, tx)
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

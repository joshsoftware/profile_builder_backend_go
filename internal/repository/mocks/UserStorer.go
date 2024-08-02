// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	pgx "github.com/jackc/pgx/v5"
	mock "github.com/stretchr/testify/mock"

	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

// UserStorer is an autogenerated mock type for the UserStorer type
type UserStorer struct {
	mock.Mock
}

// CreateUserAsEmployee provides a mock function with given fields: ctx, email, tx
func (_m *UserStorer) CreateUserAsEmployee(ctx context.Context, email string, tx pgx.Tx) error {
	ret := _m.Called(ctx, email, tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateUserAsEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, pgx.Tx) error); ok {
		r0 = rf(ctx, email, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserInfoByEmail provides a mock function with given fields: ctx, email
func (_m *UserStorer) GetUserInfoByEmail(ctx context.Context, email string) (repository.UserDao, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserInfoByEmail")
	}

	var r0 repository.UserDao
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (repository.UserDao, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) repository.UserDao); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(repository.UserDao)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveUserEmployee provides a mock function with given fields: ctx, email, tx
func (_m *UserStorer) RemoveUserEmployee(ctx context.Context, email string, tx pgx.Tx) error {
	ret := _m.Called(ctx, email, tx)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUserEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, pgx.Tx) error); ok {
		r0 = rf(ctx, email, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserStorer creates a new instance of UserStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserStorer {
	mock := &UserStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
)

// UserEmailService is an autogenerated mock type for the UserEmailService type
type UserEmailService struct {
	mock.Mock
}

// SendAdminInvitation provides a mock function with given fields: ctx, userID, request
func (_m *UserEmailService) SendAdminInvitation(ctx context.Context, userID int, request specs.UserSendInvitationRequest) error {
	ret := _m.Called(ctx, userID, request)

	if len(ret) == 0 {
		panic("no return value specified for SendAdminInvitation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.UserSendInvitationRequest) error); ok {
		r0 = rf(ctx, userID, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendUserInvitation provides a mock function with given fields: ctx, userID, request
func (_m *UserEmailService) SendUserInvitation(ctx context.Context, userID int, request specs.UserSendInvitationRequest) error {
	ret := _m.Called(ctx, userID, request)

	if len(ret) == 0 {
		panic("no return value specified for SendUserInvitation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.UserSendInvitationRequest) error); ok {
		r0 = rf(ctx, userID, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserEmailService creates a new instance of UserEmailService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserEmailService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserEmailService {
	mock := &UserEmailService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Code generated by mockery v2.40.0. DO NOT EDIT.

package mocks

import (
	context "context"

	specs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	repository "github.com/joshsoftware/profile_builder_backend_go/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// CertificateStorer is an autogenerated mock type for the CertificateStorer type
type CertificateStorer struct {
	mock.Mock
}

// CreateCertificate provides a mock function with given fields: ctx, values
func (_m *CertificateStorer) CreateCertificate(ctx context.Context, values []repository.CertificateRepo) error {
	ret := _m.Called(ctx, values)

	if len(ret) == 0 {
		panic("no return value specified for CreateCertificate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []repository.CertificateRepo) error); ok {
		r0 = rf(ctx, values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListCertificates provides a mock function with given fields: ctx, profileID, filter
func (_m *CertificateStorer) ListCertificates(ctx context.Context, profileID int, filter specs.ListCertificateFilter) ([]specs.CertificateResponse, error) {
	ret := _m.Called(ctx, profileID, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListCertificates")
	}

	var r0 []specs.CertificateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListCertificateFilter) ([]specs.CertificateResponse, error)); ok {
		return rf(ctx, profileID, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, specs.ListCertificateFilter) []specs.CertificateResponse); ok {
		r0 = rf(ctx, profileID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]specs.CertificateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, specs.ListCertificateFilter) error); ok {
		r1 = rf(ctx, profileID, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCertificate provides a mock function with given fields: ctx, profileID, eduID, req
func (_m *CertificateStorer) UpdateCertificate(ctx context.Context, profileID int, eduID int, req repository.UpdateCertificateRepo) (int, error) {
	ret := _m.Called(ctx, profileID, eduID, req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCertificate")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateCertificateRepo) (int, error)); ok {
		return rf(ctx, profileID, eduID, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, repository.UpdateCertificateRepo) int); ok {
		r0 = rf(ctx, profileID, eduID, req)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, repository.UpdateCertificateRepo) error); ok {
		r1 = rf(ctx, profileID, eduID, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCertificateStorer creates a new instance of CertificateStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCertificateStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *CertificateStorer {
	mock := &CertificateStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

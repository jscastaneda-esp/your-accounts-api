// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IUserApp is an autogenerated mock type for the IUserApp type
type IUserApp struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, uid, email
func (_m *IUserApp) Create(ctx context.Context, uid string, email string) (uint, error) {
	ret := _m.Called(ctx, uid, email)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (uint, error)); ok {
		return rf(ctx, uid, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) uint); ok {
		r0 = rf(ctx, uid, email)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, uid, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, uid, email
func (_m *IUserApp) Login(ctx context.Context, uid string, email string) (string, error) {
	ret := _m.Called(ctx, uid, email)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, uid, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, uid, email)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, uid, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIUserApp interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUserApp creates a new instance of IUserApp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUserApp(t mockConstructorTestingTNewIUserApp) *IUserApp {
	mock := &IUserApp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
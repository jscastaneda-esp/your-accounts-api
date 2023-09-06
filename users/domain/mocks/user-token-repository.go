// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "your-accounts-api/users/domain"

	mock "github.com/stretchr/testify/mock"

	persistent "your-accounts-api/shared/domain/persistent"
)

// UserTokenRepository is an autogenerated mock type for the UserTokenRepository type
type UserTokenRepository struct {
	mock.Mock
}

// Save provides a mock function with given fields: ctx, _a1
func (_m *UserTokenRepository) Save(ctx context.Context, _a1 domain.UserToken) (uint, error) {
	ret := _m.Called(ctx, _a1)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserToken) (uint, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserToken) uint); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.UserToken) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchByExample provides a mock function with given fields: ctx, example
func (_m *UserTokenRepository) SearchByExample(ctx context.Context, example domain.UserToken) (*domain.UserToken, error) {
	ret := _m.Called(ctx, example)

	var r0 *domain.UserToken
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserToken) (*domain.UserToken, error)); ok {
		return rf(ctx, example)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserToken) *domain.UserToken); ok {
		r0 = rf(ctx, example)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.UserToken) error); ok {
		r1 = rf(ctx, example)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTransaction provides a mock function with given fields: tx
func (_m *UserTokenRepository) WithTransaction(tx persistent.Transaction) domain.UserTokenRepository {
	ret := _m.Called(tx)

	var r0 domain.UserTokenRepository
	if rf, ok := ret.Get(0).(func(persistent.Transaction) domain.UserTokenRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.UserTokenRepository)
		}
	}

	return r0
}

type mockConstructorTestingTNewUserTokenRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserTokenRepository creates a new instance of UserTokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserTokenRepository(t mockConstructorTestingTNewUserTokenRepository) *UserTokenRepository {
	mock := &UserTokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
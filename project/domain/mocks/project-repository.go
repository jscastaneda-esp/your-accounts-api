// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "your-accounts-api/project/domain"

	mock "github.com/stretchr/testify/mock"

	persistent "your-accounts-api/shared/domain/persistent"
)

// ProjectRepository is an autogenerated mock type for the ProjectRepository type
type ProjectRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *ProjectRepository) Delete(ctx context.Context, id uint) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, _a1
func (_m *ProjectRepository) Save(ctx context.Context, _a1 domain.Project) (uint, error) {
	ret := _m.Called(ctx, _a1)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Project) (uint, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Project) uint); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Project) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTransaction provides a mock function with given fields: tx
func (_m *ProjectRepository) WithTransaction(tx persistent.Transaction) domain.ProjectRepository {
	ret := _m.Called(tx)

	var r0 domain.ProjectRepository
	if rf, ok := ret.Get(0).(func(persistent.Transaction) domain.ProjectRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.ProjectRepository)
		}
	}

	return r0
}

type mockConstructorTestingTNewProjectRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewProjectRepository creates a new instance of ProjectRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProjectRepository(t mockConstructorTestingTNewProjectRepository) *ProjectRepository {
	mock := &ProjectRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

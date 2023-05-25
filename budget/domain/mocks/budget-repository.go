// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	domain "api-your-accounts/budget/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"

	persistent "api-your-accounts/shared/domain/persistent"
)

// BudgetRepository is an autogenerated mock type for the BudgetRepository type
type BudgetRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, t
func (_m *BudgetRepository) Create(ctx context.Context, t domain.Budget) (*domain.Budget, error) {
	ret := _m.Called(ctx, t)

	var r0 *domain.Budget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Budget) (*domain.Budget, error)); ok {
		return rf(ctx, t)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Budget) *domain.Budget); ok {
		r0 = rf(ctx, t)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Budget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Budget) error); ok {
		r1 = rf(ctx, t)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByProjectId provides a mock function with given fields: ctx, projectId
func (_m *BudgetRepository) DeleteByProjectId(ctx context.Context, projectId uint) error {
	ret := _m.Called(ctx, projectId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, projectId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindById provides a mock function with given fields: ctx, id
func (_m *BudgetRepository) FindById(ctx context.Context, id uint) (*domain.Budget, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.Budget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*domain.Budget, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *domain.Budget); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Budget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByProjectIds provides a mock function with given fields: ctx, projectIds
func (_m *BudgetRepository) FindByProjectIds(ctx context.Context, projectIds []uint) ([]*domain.Budget, error) {
	ret := _m.Called(ctx, projectIds)

	var r0 []*domain.Budget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uint) ([]*domain.Budget, error)); ok {
		return rf(ctx, projectIds)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uint) []*domain.Budget); ok {
		r0 = rf(ctx, projectIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Budget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uint) error); ok {
		r1 = rf(ctx, projectIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTransaction provides a mock function with given fields: tx
func (_m *BudgetRepository) WithTransaction(tx persistent.Transaction) domain.BudgetRepository {
	ret := _m.Called(tx)

	var r0 domain.BudgetRepository
	if rf, ok := ret.Get(0).(func(persistent.Transaction) domain.BudgetRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.BudgetRepository)
		}
	}

	return r0
}

type mockConstructorTestingTNewBudgetRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewBudgetRepository creates a new instance of BudgetRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBudgetRepository(t mockConstructorTestingTNewBudgetRepository) *BudgetRepository {
	mock := &BudgetRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

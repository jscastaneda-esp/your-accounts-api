// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "your-accounts-api/budgets/domain"

	mock "github.com/stretchr/testify/mock"

	persistent "your-accounts-api/shared/domain/persistent"
)

// BudgetRepository is an autogenerated mock type for the BudgetRepository type
type BudgetRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *BudgetRepository) Delete(ctx context.Context, id uint) error {
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
func (_m *BudgetRepository) Save(ctx context.Context, _a1 domain.Budget) (uint, error) {
	ret := _m.Called(ctx, _a1)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Budget) (uint, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Budget) uint); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Budget) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Search provides a mock function with given fields: ctx, id
func (_m *BudgetRepository) Search(ctx context.Context, id uint) (*domain.Budget, error) {
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

// SearchAllByExample provides a mock function with given fields: ctx, example
func (_m *BudgetRepository) SearchAllByExample(ctx context.Context, example domain.Budget) ([]*domain.Budget, error) {
	ret := _m.Called(ctx, example)

	var r0 []*domain.Budget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Budget) ([]*domain.Budget, error)); ok {
		return rf(ctx, example)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Budget) []*domain.Budget); ok {
		r0 = rf(ctx, example)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Budget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Budget) error); ok {
		r1 = rf(ctx, example)
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
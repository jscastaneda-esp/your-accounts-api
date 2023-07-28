// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "your-accounts-api/budget/domain"

	mock "github.com/stretchr/testify/mock"

	persistent "your-accounts-api/shared/domain/persistent"
)

// BudgetAvailableBalanceRepository is an autogenerated mock type for the BudgetAvailableBalanceRepository type
type BudgetAvailableBalanceRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *BudgetAvailableBalanceRepository) Delete(ctx context.Context, id uint) error {
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
func (_m *BudgetAvailableBalanceRepository) Save(ctx context.Context, _a1 domain.BudgetAvailableBalance) (uint, error) {
	ret := _m.Called(ctx, _a1)

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.BudgetAvailableBalance) (uint, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.BudgetAvailableBalance) uint); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.BudgetAvailableBalance) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveAll provides a mock function with given fields: ctx, domains
func (_m *BudgetAvailableBalanceRepository) SaveAll(ctx context.Context, domains []domain.BudgetAvailableBalance) error {
	ret := _m.Called(ctx, domains)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.BudgetAvailableBalance) error); ok {
		r0 = rf(ctx, domains)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SearchAllByExample provides a mock function with given fields: ctx, example
func (_m *BudgetAvailableBalanceRepository) SearchAllByExample(ctx context.Context, example domain.BudgetAvailableBalance) ([]*domain.BudgetAvailableBalance, error) {
	ret := _m.Called(ctx, example)

	var r0 []*domain.BudgetAvailableBalance
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.BudgetAvailableBalance) ([]*domain.BudgetAvailableBalance, error)); ok {
		return rf(ctx, example)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.BudgetAvailableBalance) []*domain.BudgetAvailableBalance); ok {
		r0 = rf(ctx, example)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.BudgetAvailableBalance)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.BudgetAvailableBalance) error); ok {
		r1 = rf(ctx, example)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithTransaction provides a mock function with given fields: tx
func (_m *BudgetAvailableBalanceRepository) WithTransaction(tx persistent.Transaction) domain.BudgetAvailableBalanceRepository {
	ret := _m.Called(tx)

	var r0 domain.BudgetAvailableBalanceRepository
	if rf, ok := ret.Get(0).(func(persistent.Transaction) domain.BudgetAvailableBalanceRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.BudgetAvailableBalanceRepository)
		}
	}

	return r0
}

type mockConstructorTestingTNewBudgetAvailableBalanceRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewBudgetAvailableBalanceRepository creates a new instance of BudgetAvailableBalanceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBudgetAvailableBalanceRepository(t mockConstructorTestingTNewBudgetAvailableBalanceRepository) *BudgetAvailableBalanceRepository {
	mock := &BudgetAvailableBalanceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

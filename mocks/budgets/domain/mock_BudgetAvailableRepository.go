// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks_domain

import (
	context "context"
	domain "your-accounts-api/budgets/domain"

	mock "github.com/stretchr/testify/mock"

	persistent "your-accounts-api/shared/domain/persistent"
)

// MockBudgetAvailableRepository is an autogenerated mock type for the BudgetAvailableRepository type
type MockBudgetAvailableRepository struct {
	mock.Mock
}

type MockBudgetAvailableRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBudgetAvailableRepository) EXPECT() *MockBudgetAvailableRepository_Expecter {
	return &MockBudgetAvailableRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockBudgetAvailableRepository) Delete(ctx context.Context, id uint) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBudgetAvailableRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockBudgetAvailableRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint
func (_e *MockBudgetAvailableRepository_Expecter) Delete(ctx interface{}, id interface{}) *MockBudgetAvailableRepository_Delete_Call {
	return &MockBudgetAvailableRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockBudgetAvailableRepository_Delete_Call) Run(run func(ctx context.Context, id uint)) *MockBudgetAvailableRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *MockBudgetAvailableRepository_Delete_Call) Return(_a0 error) *MockBudgetAvailableRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBudgetAvailableRepository_Delete_Call) RunAndReturn(run func(context.Context, uint) error) *MockBudgetAvailableRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: ctx, _a1
func (_m *MockBudgetAvailableRepository) Save(ctx context.Context, _a1 domain.BudgetAvailable) (uint, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.BudgetAvailable) (uint, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.BudgetAvailable) uint); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.BudgetAvailable) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBudgetAvailableRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockBudgetAvailableRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 domain.BudgetAvailable
func (_e *MockBudgetAvailableRepository_Expecter) Save(ctx interface{}, _a1 interface{}) *MockBudgetAvailableRepository_Save_Call {
	return &MockBudgetAvailableRepository_Save_Call{Call: _e.mock.On("Save", ctx, _a1)}
}

func (_c *MockBudgetAvailableRepository_Save_Call) Run(run func(ctx context.Context, _a1 domain.BudgetAvailable)) *MockBudgetAvailableRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.BudgetAvailable))
	})
	return _c
}

func (_c *MockBudgetAvailableRepository_Save_Call) Return(_a0 uint, _a1 error) *MockBudgetAvailableRepository_Save_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBudgetAvailableRepository_Save_Call) RunAndReturn(run func(context.Context, domain.BudgetAvailable) (uint, error)) *MockBudgetAvailableRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: ctx, domains
func (_m *MockBudgetAvailableRepository) SaveAll(ctx context.Context, domains []domain.BudgetAvailable) error {
	ret := _m.Called(ctx, domains)

	if len(ret) == 0 {
		panic("no return value specified for SaveAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.BudgetAvailable) error); ok {
		r0 = rf(ctx, domains)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBudgetAvailableRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type MockBudgetAvailableRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - ctx context.Context
//   - domains []domain.BudgetAvailable
func (_e *MockBudgetAvailableRepository_Expecter) SaveAll(ctx interface{}, domains interface{}) *MockBudgetAvailableRepository_SaveAll_Call {
	return &MockBudgetAvailableRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", ctx, domains)}
}

func (_c *MockBudgetAvailableRepository_SaveAll_Call) Run(run func(ctx context.Context, domains []domain.BudgetAvailable)) *MockBudgetAvailableRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]domain.BudgetAvailable))
	})
	return _c
}

func (_c *MockBudgetAvailableRepository_SaveAll_Call) Return(_a0 error) *MockBudgetAvailableRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBudgetAvailableRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []domain.BudgetAvailable) error) *MockBudgetAvailableRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

// WithTransaction provides a mock function with given fields: tx
func (_m *MockBudgetAvailableRepository) WithTransaction(tx persistent.Transaction) domain.BudgetAvailableRepository {
	ret := _m.Called(tx)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 domain.BudgetAvailableRepository
	if rf, ok := ret.Get(0).(func(persistent.Transaction) domain.BudgetAvailableRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.BudgetAvailableRepository)
		}
	}

	return r0
}

// MockBudgetAvailableRepository_WithTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTransaction'
type MockBudgetAvailableRepository_WithTransaction_Call struct {
	*mock.Call
}

// WithTransaction is a helper method to define mock.On call
//   - tx persistent.Transaction
func (_e *MockBudgetAvailableRepository_Expecter) WithTransaction(tx interface{}) *MockBudgetAvailableRepository_WithTransaction_Call {
	return &MockBudgetAvailableRepository_WithTransaction_Call{Call: _e.mock.On("WithTransaction", tx)}
}

func (_c *MockBudgetAvailableRepository_WithTransaction_Call) Run(run func(tx persistent.Transaction)) *MockBudgetAvailableRepository_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(persistent.Transaction))
	})
	return _c
}

func (_c *MockBudgetAvailableRepository_WithTransaction_Call) Return(_a0 domain.BudgetAvailableRepository) *MockBudgetAvailableRepository_WithTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBudgetAvailableRepository_WithTransaction_Call) RunAndReturn(run func(persistent.Transaction) domain.BudgetAvailableRepository) *MockBudgetAvailableRepository_WithTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockBudgetAvailableRepository creates a new instance of MockBudgetAvailableRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBudgetAvailableRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBudgetAvailableRepository {
	mock := &MockBudgetAvailableRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
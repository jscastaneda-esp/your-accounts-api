// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks_application

import (
	context "context"
	application "your-accounts-api/budgets/application"

	domain "your-accounts-api/budgets/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockIBudgetApp is an autogenerated mock type for the IBudgetApp type
type MockIBudgetApp struct {
	mock.Mock
}

type MockIBudgetApp_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIBudgetApp) EXPECT() *MockIBudgetApp_Expecter {
	return &MockIBudgetApp_Expecter{mock: &_m.Mock}
}

// Changes provides a mock function with given fields: ctx, id, changes
func (_m *MockIBudgetApp) Changes(ctx context.Context, id uint, changes []application.Change) []application.ChangeResult {
	ret := _m.Called(ctx, id, changes)

	if len(ret) == 0 {
		panic("no return value specified for Changes")
	}

	var r0 []application.ChangeResult
	if rf, ok := ret.Get(0).(func(context.Context, uint, []application.Change) []application.ChangeResult); ok {
		r0 = rf(ctx, id, changes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]application.ChangeResult)
		}
	}

	return r0
}

// MockIBudgetApp_Changes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Changes'
type MockIBudgetApp_Changes_Call struct {
	*mock.Call
}

// Changes is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint
//   - changes []application.Change
func (_e *MockIBudgetApp_Expecter) Changes(ctx interface{}, id interface{}, changes interface{}) *MockIBudgetApp_Changes_Call {
	return &MockIBudgetApp_Changes_Call{Call: _e.mock.On("Changes", ctx, id, changes)}
}

func (_c *MockIBudgetApp_Changes_Call) Run(run func(ctx context.Context, id uint, changes []application.Change)) *MockIBudgetApp_Changes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].([]application.Change))
	})
	return _c
}

func (_c *MockIBudgetApp_Changes_Call) Return(_a0 []application.ChangeResult) *MockIBudgetApp_Changes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIBudgetApp_Changes_Call) RunAndReturn(run func(context.Context, uint, []application.Change) []application.ChangeResult) *MockIBudgetApp_Changes_Call {
	_c.Call.Return(run)
	return _c
}

// Clone provides a mock function with given fields: ctx, userId, baseId
func (_m *MockIBudgetApp) Clone(ctx context.Context, userId uint, baseId uint) (uint, error) {
	ret := _m.Called(ctx, userId, baseId)

	if len(ret) == 0 {
		panic("no return value specified for Clone")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) (uint, error)); ok {
		return rf(ctx, userId, baseId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) uint); ok {
		r0 = rf(ctx, userId, baseId)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, userId, baseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIBudgetApp_Clone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clone'
type MockIBudgetApp_Clone_Call struct {
	*mock.Call
}

// Clone is a helper method to define mock.On call
//   - ctx context.Context
//   - userId uint
//   - baseId uint
func (_e *MockIBudgetApp_Expecter) Clone(ctx interface{}, userId interface{}, baseId interface{}) *MockIBudgetApp_Clone_Call {
	return &MockIBudgetApp_Clone_Call{Call: _e.mock.On("Clone", ctx, userId, baseId)}
}

func (_c *MockIBudgetApp_Clone_Call) Run(run func(ctx context.Context, userId uint, baseId uint)) *MockIBudgetApp_Clone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(uint))
	})
	return _c
}

func (_c *MockIBudgetApp_Clone_Call) Return(_a0 uint, _a1 error) *MockIBudgetApp_Clone_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIBudgetApp_Clone_Call) RunAndReturn(run func(context.Context, uint, uint) (uint, error)) *MockIBudgetApp_Clone_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, userId, name
func (_m *MockIBudgetApp) Create(ctx context.Context, userId uint, name string) (uint, error) {
	ret := _m.Called(ctx, userId, name)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) (uint, error)); ok {
		return rf(ctx, userId, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) uint); ok {
		r0 = rf(ctx, userId, name)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, string) error); ok {
		r1 = rf(ctx, userId, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIBudgetApp_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockIBudgetApp_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - userId uint
//   - name string
func (_e *MockIBudgetApp_Expecter) Create(ctx interface{}, userId interface{}, name interface{}) *MockIBudgetApp_Create_Call {
	return &MockIBudgetApp_Create_Call{Call: _e.mock.On("Create", ctx, userId, name)}
}

func (_c *MockIBudgetApp_Create_Call) Run(run func(ctx context.Context, userId uint, name string)) *MockIBudgetApp_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(string))
	})
	return _c
}

func (_c *MockIBudgetApp_Create_Call) Return(_a0 uint, _a1 error) *MockIBudgetApp_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIBudgetApp_Create_Call) RunAndReturn(run func(context.Context, uint, string) (uint, error)) *MockIBudgetApp_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockIBudgetApp) Delete(ctx context.Context, id uint) error {
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

// MockIBudgetApp_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockIBudgetApp_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint
func (_e *MockIBudgetApp_Expecter) Delete(ctx interface{}, id interface{}) *MockIBudgetApp_Delete_Call {
	return &MockIBudgetApp_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockIBudgetApp_Delete_Call) Run(run func(ctx context.Context, id uint)) *MockIBudgetApp_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *MockIBudgetApp_Delete_Call) Return(_a0 error) *MockIBudgetApp_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIBudgetApp_Delete_Call) RunAndReturn(run func(context.Context, uint) error) *MockIBudgetApp_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindById provides a mock function with given fields: ctx, id
func (_m *MockIBudgetApp) FindById(ctx context.Context, id uint) (domain.Budget, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindById")
	}

	var r0 domain.Budget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (domain.Budget, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) domain.Budget); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Budget)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIBudgetApp_FindById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindById'
type MockIBudgetApp_FindById_Call struct {
	*mock.Call
}

// FindById is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint
func (_e *MockIBudgetApp_Expecter) FindById(ctx interface{}, id interface{}) *MockIBudgetApp_FindById_Call {
	return &MockIBudgetApp_FindById_Call{Call: _e.mock.On("FindById", ctx, id)}
}

func (_c *MockIBudgetApp_FindById_Call) Run(run func(ctx context.Context, id uint)) *MockIBudgetApp_FindById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *MockIBudgetApp_FindById_Call) Return(_a0 domain.Budget, _a1 error) *MockIBudgetApp_FindById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIBudgetApp_FindById_Call) RunAndReturn(run func(context.Context, uint) (domain.Budget, error)) *MockIBudgetApp_FindById_Call {
	_c.Call.Return(run)
	return _c
}

// FindByUserId provides a mock function with given fields: ctx, userId
func (_m *MockIBudgetApp) FindByUserId(ctx context.Context, userId uint) ([]domain.Budget, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindByUserId")
	}

	var r0 []domain.Budget
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) ([]domain.Budget, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) []domain.Budget); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Budget)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIBudgetApp_FindByUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByUserId'
type MockIBudgetApp_FindByUserId_Call struct {
	*mock.Call
}

// FindByUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - userId uint
func (_e *MockIBudgetApp_Expecter) FindByUserId(ctx interface{}, userId interface{}) *MockIBudgetApp_FindByUserId_Call {
	return &MockIBudgetApp_FindByUserId_Call{Call: _e.mock.On("FindByUserId", ctx, userId)}
}

func (_c *MockIBudgetApp_FindByUserId_Call) Run(run func(ctx context.Context, userId uint)) *MockIBudgetApp_FindByUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *MockIBudgetApp_FindByUserId_Call) Return(_a0 []domain.Budget, _a1 error) *MockIBudgetApp_FindByUserId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIBudgetApp_FindByUserId_Call) RunAndReturn(run func(context.Context, uint) ([]domain.Budget, error)) *MockIBudgetApp_FindByUserId_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIBudgetApp creates a new instance of MockIBudgetApp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIBudgetApp(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIBudgetApp {
	mock := &MockIBudgetApp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

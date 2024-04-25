// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks_application

import (
	context "context"
	domain "your-accounts-api/shared/domain"

	mock "github.com/stretchr/testify/mock"

	persistent "your-accounts-api/shared/domain/persistent"
)

// MockILogApp is an autogenerated mock type for the ILogApp type
type MockILogApp struct {
	mock.Mock
}

type MockILogApp_Expecter struct {
	mock *mock.Mock
}

func (_m *MockILogApp) EXPECT() *MockILogApp_Expecter {
	return &MockILogApp_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, description, code, resourceId, detail, tx
func (_m *MockILogApp) Create(ctx context.Context, description string, code domain.CodeLog, resourceId uint, detail map[string]interface{}, tx persistent.Transaction) error {
	ret := _m.Called(ctx, description, code, resourceId, detail, tx)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.CodeLog, uint, map[string]interface{}, persistent.Transaction) error); ok {
		r0 = rf(ctx, description, code, resourceId, detail, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockILogApp_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockILogApp_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - description string
//   - code domain.CodeLog
//   - resourceId uint
//   - detail map[string]interface{}
//   - tx persistent.Transaction
func (_e *MockILogApp_Expecter) Create(ctx interface{}, description interface{}, code interface{}, resourceId interface{}, detail interface{}, tx interface{}) *MockILogApp_Create_Call {
	return &MockILogApp_Create_Call{Call: _e.mock.On("Create", ctx, description, code, resourceId, detail, tx)}
}

func (_c *MockILogApp_Create_Call) Run(run func(ctx context.Context, description string, code domain.CodeLog, resourceId uint, detail map[string]interface{}, tx persistent.Transaction)) *MockILogApp_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(domain.CodeLog), args[3].(uint), args[4].(map[string]interface{}), args[5].(persistent.Transaction))
	})
	return _c
}

func (_c *MockILogApp_Create_Call) Return(_a0 error) *MockILogApp_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockILogApp_Create_Call) RunAndReturn(run func(context.Context, string, domain.CodeLog, uint, map[string]interface{}, persistent.Transaction) error) *MockILogApp_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOld provides a mock function with given fields: ctx
func (_m *MockILogApp) DeleteOld(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOld")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockILogApp_DeleteOld_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOld'
type MockILogApp_DeleteOld_Call struct {
	*mock.Call
}

// DeleteOld is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockILogApp_Expecter) DeleteOld(ctx interface{}) *MockILogApp_DeleteOld_Call {
	return &MockILogApp_DeleteOld_Call{Call: _e.mock.On("DeleteOld", ctx)}
}

func (_c *MockILogApp_DeleteOld_Call) Run(run func(ctx context.Context)) *MockILogApp_DeleteOld_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockILogApp_DeleteOld_Call) Return(_a0 error) *MockILogApp_DeleteOld_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockILogApp_DeleteOld_Call) RunAndReturn(run func(context.Context) error) *MockILogApp_DeleteOld_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOrphan provides a mock function with given fields: ctx
func (_m *MockILogApp) DeleteOrphan(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrphan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockILogApp_DeleteOrphan_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOrphan'
type MockILogApp_DeleteOrphan_Call struct {
	*mock.Call
}

// DeleteOrphan is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockILogApp_Expecter) DeleteOrphan(ctx interface{}) *MockILogApp_DeleteOrphan_Call {
	return &MockILogApp_DeleteOrphan_Call{Call: _e.mock.On("DeleteOrphan", ctx)}
}

func (_c *MockILogApp_DeleteOrphan_Call) Run(run func(ctx context.Context)) *MockILogApp_DeleteOrphan_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockILogApp_DeleteOrphan_Call) Return(_a0 error) *MockILogApp_DeleteOrphan_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockILogApp_DeleteOrphan_Call) RunAndReturn(run func(context.Context) error) *MockILogApp_DeleteOrphan_Call {
	_c.Call.Return(run)
	return _c
}

// FindByProject provides a mock function with given fields: ctx, code, resourceId
func (_m *MockILogApp) FindByProject(ctx context.Context, code domain.CodeLog, resourceId uint) ([]domain.Log, error) {
	ret := _m.Called(ctx, code, resourceId)

	if len(ret) == 0 {
		panic("no return value specified for FindByProject")
	}

	var r0 []domain.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CodeLog, uint) ([]domain.Log, error)); ok {
		return rf(ctx, code, resourceId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.CodeLog, uint) []domain.Log); ok {
		r0 = rf(ctx, code, resourceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.CodeLog, uint) error); ok {
		r1 = rf(ctx, code, resourceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockILogApp_FindByProject_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByProject'
type MockILogApp_FindByProject_Call struct {
	*mock.Call
}

// FindByProject is a helper method to define mock.On call
//   - ctx context.Context
//   - code domain.CodeLog
//   - resourceId uint
func (_e *MockILogApp_Expecter) FindByProject(ctx interface{}, code interface{}, resourceId interface{}) *MockILogApp_FindByProject_Call {
	return &MockILogApp_FindByProject_Call{Call: _e.mock.On("FindByProject", ctx, code, resourceId)}
}

func (_c *MockILogApp_FindByProject_Call) Run(run func(ctx context.Context, code domain.CodeLog, resourceId uint)) *MockILogApp_FindByProject_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.CodeLog), args[2].(uint))
	})
	return _c
}

func (_c *MockILogApp_FindByProject_Call) Return(_a0 []domain.Log, _a1 error) *MockILogApp_FindByProject_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockILogApp_FindByProject_Call) RunAndReturn(run func(context.Context, domain.CodeLog, uint) ([]domain.Log, error)) *MockILogApp_FindByProject_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockILogApp creates a new instance of MockILogApp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockILogApp(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockILogApp {
	mock := &MockILogApp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
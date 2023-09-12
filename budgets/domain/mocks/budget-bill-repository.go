// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// BudgetBillRepository is an autogenerated mock type for the BudgetBillRepository type
type BudgetBillRepository struct {
	mock.Mock
}

type mockConstructorTestingTNewBudgetBillRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewBudgetBillRepository creates a new instance of BudgetBillRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBudgetBillRepository(t mockConstructorTestingTNewBudgetBillRepository) *BudgetBillRepository {
	mock := &BudgetBillRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

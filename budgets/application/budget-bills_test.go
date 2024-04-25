package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/budgets/domain"
	mocks_domain "your-accounts-api/mocks/budgets/domain"
	mocks_application "your-accounts-api/mocks/shared/application"
	mocks_persistent "your-accounts-api/mocks/shared/domain/persistent"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestBudgetBillSuite struct {
	suite.Suite
	id                     uint
	description            string
	category               domain.BudgetBillCategory
	budgetId               uint
	mockTransactionManager *mocks_persistent.MockTransactionManager
	mockBudgetBillRepo     *mocks_domain.MockBudgetBillRepository
	mockLogApp             *mocks_application.MockILogApp
	app                    IBudgetBillApp
	ctx                    context.Context
}

func (suite *TestBudgetBillSuite) SetupSuite() {
	suite.id = 1
	suite.description = "Test"
	suite.category = domain.Education
	suite.budgetId = 1
	suite.ctx = context.Background()
}

func (suite *TestBudgetBillSuite) SetupTest() {
	suite.mockTransactionManager = mocks_persistent.NewMockTransactionManager(suite.T())
	suite.mockBudgetBillRepo = mocks_domain.NewMockBudgetBillRepository(suite.T())
	suite.mockLogApp = mocks_application.NewMockILogApp(suite.T())
	suite.app = NewBudgetBillApp(suite.mockTransactionManager, suite.mockBudgetBillRepo, suite.mockLogApp)
}

func (suite *TestBudgetBillSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Create(suite.ctx, "Test", domain.Education, suite.budgetId)

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetBillSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project log")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", domain.Education, suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetBillSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", domain.Education, suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetBillSuite) TestCreateError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation bill")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", domain.Education, suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetBillSuite) TestCreateTransactionSuccess() {
	require := require.New(suite.T())
	payment := float64(0)
	billExpected := domain.BudgetBill{
		ID:          &suite.id,
		Description: &suite.description,
		Payment:     &payment,
		Category:    &suite.category,
		BudgetId:    &suite.budgetId,
	}
	suite.mockBudgetBillRepo.On("Search", suite.ctx, suite.id).Return(billExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.BudgetBill, suite.id, mock.Anything, nil).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("Save", suite.ctx, mock.Anything).Return(suite.id, nil)

	err := suite.app.CreateTransaction(suite.ctx, "Trans 1", float64(10000), suite.id)

	require.NoError(err)
}

func (suite *TestBudgetBillSuite) TestCreateTransactionErrorSearch() {
	require := require.New(suite.T())
	suite.mockBudgetBillRepo.On("Search", suite.ctx, suite.id).Return(domain.BudgetBill{}, gorm.ErrRecordNotFound)

	err := suite.app.CreateTransaction(suite.ctx, "Trans 1", float64(10000), suite.id)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
}

func (suite *TestBudgetBillSuite) TestCreateTransactionErrorCreateLog() {
	require := require.New(suite.T())
	payment := float64(0)
	billExpected := domain.BudgetBill{
		ID:          &suite.id,
		Description: &suite.description,
		Payment:     &payment,
		Category:    &suite.category,
		BudgetId:    &suite.budgetId,
	}
	errExpected := errors.New("Error in creation project log")
	suite.mockBudgetBillRepo.On("Search", suite.ctx, suite.id).Return(billExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.BudgetBill, suite.id, mock.Anything, nil).Return(errExpected)

	err := suite.app.CreateTransaction(suite.ctx, "Trans 1", float64(10000), suite.id)

	require.EqualError(errExpected, err.Error())
}

func TestTestBudgetBillSuite(t *testing.T) {
	suite.Run(t, new(TestBudgetBillSuite))
}

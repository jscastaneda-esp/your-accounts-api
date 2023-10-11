package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/domain/mocks"
	mocks_logs "your-accounts-api/shared/application/mocks"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestBudgetBillSuite struct {
	suite.Suite
	budgetId               uint
	mockTransactionManager *mocks_shared.TransactionManager
	mockBudgetBillRepo     *mocks.BudgetBillRepository
	mockLogApp             *mocks_logs.ILogApp
	app                    IBudgetBillApp
	ctx                    context.Context
}

func (suite *TestBudgetBillSuite) SetupSuite() {
	suite.budgetId = 1
	suite.ctx = context.Background()
}

func (suite *TestBudgetBillSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockBudgetBillRepo = mocks.NewBudgetBillRepository(suite.T())
	suite.mockLogApp = mocks_logs.NewILogApp(suite.T())
	suite.app = NewBudgetBillApp(suite.mockTransactionManager, suite.mockBudgetBillRepo, suite.mockLogApp)
}

func (suite *TestBudgetBillSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("CreateLog", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)
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
	suite.mockLogApp.On("CreateLog", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(errExpected)

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
	errExpected := errors.New("Error in creation available")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("CreateLog", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", domain.Education, suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func TestTestBudgetBillSuite(t *testing.T) {
	suite.Run(t, new(TestBudgetBillSuite))
}

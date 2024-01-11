package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/budgets/domain/mocks"
	mocks_logs "your-accounts-api/shared/application/mocks"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestBudgetAvailableSuite struct {
	suite.Suite
	budgetId                uint
	mockTransactionManager  *mocks_shared.TransactionManager
	mockBudgetAvailableRepo *mocks.BudgetAvailableRepository
	mockLogApp              *mocks_logs.ILogApp
	app                     IBudgetAvailableApp
	ctx                     context.Context
}

func (suite *TestBudgetAvailableSuite) SetupSuite() {
	suite.budgetId = 1
	suite.ctx = context.Background()
}

func (suite *TestBudgetAvailableSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockBudgetAvailableRepo = mocks.NewBudgetAvailableRepository(suite.T())
	suite.mockLogApp = mocks_logs.NewILogApp(suite.T())
	suite.app = NewBudgetAvailableApp(suite.mockTransactionManager, suite.mockBudgetAvailableRepo, suite.mockLogApp)
}

func (suite *TestBudgetAvailableSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo)
	suite.mockBudgetAvailableRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetAvailableSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project log")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetAvailableSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetAvailableSuite) TestCreateError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation available")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo)
	suite.mockBudgetAvailableRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func TestTestBudgetAvailableSuite(t *testing.T) {
	suite.Run(t, new(TestBudgetAvailableSuite))
}

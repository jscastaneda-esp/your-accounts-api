package application

import (
	"context"
	"errors"
	"testing"
	mocks_budgets "your-accounts-api/budgets/application/mocks"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/domain/mocks"
	mocks_logs "your-accounts-api/shared/application/mocks"
	"your-accounts-api/shared/domain/persistent"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestBudgetAvailableSuite struct {
	suite.Suite
	budgetId                       uint
	mockTransactionManager         *mocks_shared.TransactionManager
	mockBudgetAvailableBalanceRepo *mocks.BudgetAvailableBalanceRepository
	mockBudgetApp                  *mocks_budgets.IBudgetApp
	mockLogApp                     *mocks_logs.ILogApp
	app                            IBudgetAvailableBalanceApp
	ctx                            context.Context
}

func (suite *TestBudgetAvailableSuite) SetupSuite() {
	suite.budgetId = 1
	suite.ctx = context.Background()
}

func (suite *TestBudgetAvailableSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockBudgetAvailableBalanceRepo = mocks.NewBudgetAvailableBalanceRepository(suite.T())
	suite.mockBudgetApp = mocks_budgets.NewIBudgetApp(suite.T())
	suite.mockLogApp = mocks_logs.NewILogApp(suite.T())
	suite.app = NewBudgetAvailableBalanceApp(suite.mockTransactionManager, suite.mockBudgetAvailableBalanceRepo, suite.mockBudgetApp, suite.mockLogApp)
}

func (suite *TestBudgetAvailableSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	userId := uint(1)
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &userId,
	}
	suite.mockBudgetApp.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("CreateLog", suite.ctx, mock.Anything, *budgetExpected.ID, mock.Anything, nil).Return(nil)
	suite.mockBudgetAvailableBalanceRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableBalanceRepo)
	suite.mockBudgetAvailableBalanceRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetAvailableSuite) TestCreateErrorFindById() {
	require := require.New(suite.T())
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetApp.On("FindById", suite.ctx, suite.budgetId).Return(nil, errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetAvailableSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	userId := uint(1)
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &userId,
	}
	errExpected := errors.New("Error in creation project log")
	suite.mockBudgetApp.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("CreateLog", suite.ctx, mock.Anything, *budgetExpected.ID, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetAvailableSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	userId := uint(1)
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &userId,
	}
	errExpected := errors.New("Error in transaction")
	suite.mockBudgetApp.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetAvailableSuite) TestCreateError() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	userId := uint(1)
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &userId,
	}
	errExpected := errors.New("Error in creation available")
	suite.mockBudgetApp.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogApp.On("CreateLog", suite.ctx, mock.Anything, *budgetExpected.ID, mock.Anything, nil).Return(nil)
	suite.mockBudgetAvailableBalanceRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableBalanceRepo)
	suite.mockBudgetAvailableBalanceRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, "Test", suite.budgetId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func TestTestBudgetAvailableSuite(t *testing.T) {
	suite.Run(t, new(TestBudgetAvailableSuite))
}

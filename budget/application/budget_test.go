package application

import (
	"api-your-accounts/budget/domain"
	"api-your-accounts/budget/domain/mocks"
	mocksShared "api-your-accounts/shared/domain/persistent/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	budgetId               uint
	mockTransactionManager *mocksShared.TransactionManager
	mockBudgetRepo         *mocks.BudgetRepository
	app                    IBudgetApp
}

func (suite *TestSuite) SetupSuite() {
	suite.budgetId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mockTransactionManager = mocksShared.NewTransactionManager(suite.T())
	suite.mockBudgetRepo = mocks.NewBudgetRepository(suite.T())
	suite.app = NewBudgetApp(suite.mockTransactionManager, suite.mockBudgetRepo)
}

func (suite *TestSuite) TestFindByIdSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	budgetExpected := &domain.Budget{
		ID:        suite.budgetId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	suite.mockBudgetRepo.On("FindById", ctx, suite.budgetId).Return(budgetExpected, nil)

	res, err := suite.app.FindById(ctx, suite.budgetId)

	require.NoError(err)
	require.Equal(budgetExpected, res)
}

func (suite *TestSuite) TestFindByIdError() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockBudgetRepo.On("FindById", ctx, suite.budgetId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindById(ctx, suite.budgetId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

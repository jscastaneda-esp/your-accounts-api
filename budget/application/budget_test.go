package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/domain/mocks"
	mocksProject "your-accounts-api/project/application/mocks"
	projectDom "your-accounts-api/project/domain"
	"your-accounts-api/shared/domain/persistent"
	mocksShared "your-accounts-api/shared/domain/persistent/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	budgetId               uint
	userId                 uint
	projectId              uint
	mockTransactionManager *mocksShared.TransactionManager
	mockBudgetRepo         *mocks.BudgetRepository
	mockProjectApp         *mocksProject.IProjectApp
	app                    IBudgetApp
	ctx                    context.Context
}

func (suite *TestSuite) SetupSuite() {
	suite.budgetId = 1
	suite.userId = 2
	suite.projectId = 3
	suite.ctx = context.Background()
}

func (suite *TestSuite) SetupTest() {
	suite.mockTransactionManager = mocksShared.NewTransactionManager(suite.T())
	suite.mockBudgetRepo = mocks.NewBudgetRepository(suite.T())
	suite.mockProjectApp = mocksProject.NewIProjectApp(suite.T())
	instance = nil
	suite.app = NewBudgetApp(suite.mockTransactionManager, suite.mockBudgetRepo, suite.mockProjectApp)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestSuite) TestCreateErrorCreateProject() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project log")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneSuccess() {
	require := require.New(suite.T())
	baseId := uint(999)
	budgetExpected := &domain.Budget{
		ID:        baseId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	suite.mockBudgetRepo.On("FindById", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestSuite) TestCloneErrorFindById() {
	require := require.New(suite.T())
	baseId := uint(999)
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("FindById", suite.ctx, baseId).Return(nil, errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorCreateProject() {
	require := require.New(suite.T())
	baseId := uint(999)
	budgetExpected := &domain.Budget{
		ID:        baseId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockBudgetRepo.On("FindById", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(uint(0), errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorCreateLog() {
	require := require.New(suite.T())
	baseId := uint(999)
	budgetExpected := &domain.Budget{
		ID:        baseId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	errExpected := errors.New("Error in creation project log")
	suite.mockBudgetRepo.On("FindById", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorTransaction() {
	require := require.New(suite.T())
	baseId := uint(999)
	budgetExpected := &domain.Budget{
		ID:        baseId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	errExpected := errors.New("Error in transaction")
	suite.mockBudgetRepo.On("FindById", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneError() {
	require := require.New(suite.T())
	baseId := uint(999)
	budgetExpected := &domain.Budget{
		ID:        baseId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	errExpected := errors.New("Error in creation budget")
	suite.mockBudgetRepo.On("FindById", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, projectDom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestFindByIdSuccess() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        suite.budgetId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	suite.mockBudgetRepo.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)

	res, err := suite.app.FindById(suite.ctx, suite.budgetId)

	require.NoError(err)
	require.Equal(budgetExpected, res)
}

func (suite *TestSuite) TestFindByIdError() {
	require := require.New(suite.T())
	suite.mockBudgetRepo.On("FindById", suite.ctx, suite.budgetId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindById(suite.ctx, suite.budgetId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestSuite) TestFindByUserIdSuccess() {
	require := require.New(suite.T())
	projectIds := []uint{999, 1000}
	budgetsExpected := []*domain.Budget{
		{
			ID:        999,
			Name:      "Test 1",
			Year:      2023,
			Month:     5,
			ProjectId: projectIds[0],
		},
		{
			ID:        1000,
			Name:      "Test 2",
			Year:      2023,
			Month:     5,
			ProjectId: projectIds[1],
		},
	}
	suite.mockProjectApp.On("FindByUserIdAndType", suite.ctx, suite.userId, projectDom.TypeBudget).Return(projectIds, nil)
	suite.mockBudgetRepo.On("FindByProjectIds", suite.ctx, projectIds).Return(budgetsExpected, nil)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.NoError(err)
	require.Equal(budgetsExpected, res)
}

func (suite *TestSuite) TestFindByUserIdErrorFindProjectsByUserIdAndType() {
	require := require.New(suite.T())
	suite.mockProjectApp.On("FindByUserIdAndType", suite.ctx, suite.userId, projectDom.TypeBudget).Return(nil, gorm.ErrInvalidField)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestFindByUserIdErrorFindByProjectIds() {
	require := require.New(suite.T())
	projectIds := []uint{999, 1000}
	suite.mockProjectApp.On("FindByUserIdAndType", suite.ctx, suite.userId, projectDom.TypeBudget).Return(projectIds, nil)
	suite.mockBudgetRepo.On("FindByProjectIds", suite.ctx, projectIds).Return(nil, gorm.ErrInvalidField)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        suite.budgetId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	suite.mockBudgetRepo.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("Delete", suite.ctx, budgetExpected.ID).Return(nil)
	suite.mockProjectApp.On("Delete", suite.ctx, budgetExpected.ProjectId, nil).Return(nil)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteErrorFindById() {
	require := require.New(suite.T())
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("FindById", suite.ctx, suite.budgetId).Return(nil, errExpected)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.EqualError(errExpected, err.Error())
}

func (suite *TestSuite) TestDeleteError() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        suite.budgetId,
		Name:      "Test",
		Year:      1,
		Month:     1,
		ProjectId: 1,
	}
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("FindById", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("Delete", suite.ctx, budgetExpected.ID).Return(errExpected)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.EqualError(errExpected, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

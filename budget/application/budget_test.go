package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/budget/domain"
	"your-accounts-api/budget/domain/mocks"
	mocks_project "your-accounts-api/project/application/mocks"
	project_dom "your-accounts-api/project/domain"
	"your-accounts-api/shared/domain/persistent"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestBudgetSuite struct {
	suite.Suite
	budgetId               uint
	userId                 uint
	projectId              uint
	mockTransactionManager *mocks_shared.TransactionManager
	mockBudgetRepo         *mocks.BudgetRepository
	mockProjectApp         *mocks_project.IProjectApp
	app                    IBudgetApp
	ctx                    context.Context
}

func (suite *TestBudgetSuite) SetupSuite() {
	suite.budgetId = 1
	suite.userId = 2
	suite.projectId = 3
	suite.ctx = context.Background()
}

func (suite *TestBudgetSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockBudgetRepo = mocks.NewBudgetRepository(suite.T())
	suite.mockProjectApp = mocks_project.NewIProjectApp(suite.T())
	suite.app = NewBudgetApp(suite.mockTransactionManager, suite.mockBudgetRepo, suite.mockProjectApp)
}

func (suite *TestBudgetSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetSuite) TestCreateErrorCreateProject() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project log")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCreateError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation budget")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneSuccess() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &baseId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetSuite) TestCloneErrorSearch() {
	require := require.New(suite.T())
	baseId := uint(999)
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(nil, errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorCreateProject() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &baseId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(uint(0), errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorCreateLog() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &baseId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	errExpected := errors.New("Error in creation project log")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorTransaction() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &baseId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	errExpected := errors.New("Error in transaction")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneError() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &baseId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	errExpected := errors.New("Error in creation budget")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectApp.On("Create", suite.ctx, suite.userId, project_dom.TypeBudget, nil).Return(suite.projectId, nil)
	suite.mockProjectApp.On("CreateLog", suite.ctx, mock.Anything, suite.projectId, mock.Anything, nil).Return(nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestFindByIdSuccess() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &suite.budgetId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)

	res, err := suite.app.FindById(suite.ctx, suite.budgetId)

	require.NoError(err)
	require.Equal(budgetExpected, res)
}

func (suite *TestBudgetSuite) TestFindByIdError() {
	require := require.New(suite.T())
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindById(suite.ctx, suite.budgetId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestBudgetSuite) TestFindByUserIdSuccess() {
	require := require.New(suite.T())
	ids := []uint{uint(999), uint(1000)}
	names := []string{"Test 1", "Test 2"}
	year := uint16(2003)
	month := uint8(5)
	projectIds := []uint{999, 1000}
	budgetsExpected := []*domain.Budget{
		{
			ID:        &ids[0],
			Name:      &names[0],
			Year:      &year,
			Month:     &month,
			ProjectId: &projectIds[0],
		},
		{
			ID:        &ids[1],
			Name:      &names[1],
			Year:      &year,
			Month:     &month,
			ProjectId: &projectIds[1],
		},
	}
	suite.mockBudgetRepo.On("SearchByUserId", suite.ctx, suite.userId).Return(budgetsExpected, nil)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.NoError(err)
	require.Equal(budgetsExpected, res)
}

func (suite *TestBudgetSuite) TestFindByUserIdError() {
	require := require.New(suite.T())
	suite.mockBudgetRepo.On("SearchByUserId", suite.ctx, suite.userId).Return(nil, gorm.ErrInvalidField)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestBudgetSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &suite.budgetId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("Delete", suite.ctx, *budgetExpected.ID).Return(nil)
	suite.mockProjectApp.On("Delete", suite.ctx, *budgetExpected.ProjectId, nil).Return(nil)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.NoError(err)
}

func (suite *TestBudgetSuite) TestDeleteErrorSearch() {
	require := require.New(suite.T())
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(nil, errExpected)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.EqualError(errExpected, err.Error())
}

func (suite *TestBudgetSuite) TestDeleteError() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:        &suite.budgetId,
		Name:      &name,
		Year:      &year,
		Month:     &month,
		ProjectId: &suite.projectId,
	}
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("Delete", suite.ctx, *budgetExpected.ID).Return(errExpected)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.EqualError(errExpected, err.Error())
}

func TestTestBudgetSuite(t *testing.T) {
	suite.Run(t, new(TestBudgetSuite))
}

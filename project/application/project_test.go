package application

import (
	"context"
	"errors"
	"testing"
	budgetDom "your-accounts-api/budget/domain"
	mocksBudget "your-accounts-api/budget/domain/mocks"
	"your-accounts-api/project/domain"
	"your-accounts-api/project/domain/mocks"
	"your-accounts-api/shared/domain/persistent"
	mocksShared "your-accounts-api/shared/domain/persistent/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	name                   string
	userId                 uint
	typeBudget             domain.ProjectType
	cloneId                uint
	mockTransactionManager *mocksShared.TransactionManager
	mockProjectRepo        *mocks.ProjectRepository
	mockProjectLogRepo     *mocks.ProjectLogRepository
	mockBudgetRepo         *mocksBudget.BudgetRepository
	app                    IProjectApp
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.userId = 1
	suite.typeBudget = domain.TypeBudget
	suite.cloneId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mockTransactionManager = mocksShared.NewTransactionManager(suite.T())
	suite.mockProjectRepo = mocks.NewProjectRepository(suite.T())
	suite.mockProjectLogRepo = mocks.NewProjectLogRepository(suite.T())
	suite.mockBudgetRepo = mocksBudget.NewBudgetRepository(suite.T())
	suite.app = NewProjectApp(suite.mockTransactionManager, suite.mockProjectRepo, suite.mockProjectLogRepo, suite.mockBudgetRepo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	createData := CreateData{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(999), nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", ctx, mock.Anything).Return(uint(0), nil)
	suite.mockProjectLogRepo.On("WithTransaction", nil).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Create", ctx, mock.Anything).Return(uint(0), nil)

	res, err := suite.app.Create(ctx, createData)

	require.NoError(err)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	ctx := context.Background()
	createData := CreateData{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(ctx, createData)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorCreateProject() {
	require := require.New(suite.T())
	ctx := context.Background()
	createData := CreateData{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(ctx, createData)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorCreateBudget() {
	require := require.New(suite.T())
	ctx := context.Background()
	createData := CreateData{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(1000), nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(ctx, createData)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	ctx := context.Background()
	createData := CreateData{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(999), nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", ctx, mock.Anything).Return(uint(0), nil)
	suite.mockProjectLogRepo.On("WithTransaction", nil).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Create", ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(ctx, createData)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	baseProject := &domain.Project{
		ID:     suite.cloneId,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	budgetsExpected := []*budgetDom.Budget{
		{
			ID:        999,
			Name:      "Test 1",
			Year:      2023,
			Month:     5,
			ProjectId: baseProject.ID,
		},
	}
	suite.mockProjectRepo.On("FindById", ctx, baseProject.ID).Return(baseProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(1000), nil)
	suite.mockBudgetRepo.On("FindByProjectIds", ctx, []uint{baseProject.ID}).Return(budgetsExpected, nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", ctx, mock.Anything).Return(uint(0), nil)
	suite.mockProjectLogRepo.On("WithTransaction", nil).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Create", ctx, mock.Anything).Return(uint(0), nil)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.NoError(err)
	require.Equal(uint(1000), res)
}

func (suite *TestSuite) TestCloneErrorFind() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockProjectRepo.On("FindById", ctx, suite.cloneId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorCreateProject() {
	require := require.New(suite.T())
	ctx := context.Background()
	baseProject := &domain.Project{
		ID:     suite.cloneId,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	suite.mockProjectRepo.On("FindById", ctx, baseProject.ID).Return(baseProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(0), gorm.ErrInvalidField)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorFindBudget() {
	require := require.New(suite.T())
	ctx := context.Background()
	baseProject := &domain.Project{
		ID:     suite.cloneId,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	suite.mockProjectRepo.On("FindById", ctx, baseProject.ID).Return(baseProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(1000), nil)
	suite.mockBudgetRepo.On("FindByProjectIds", ctx, []uint{baseProject.ID}).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorCreateBudget() {
	require := require.New(suite.T())
	ctx := context.Background()
	baseProject := &domain.Project{
		ID:     suite.cloneId,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	budgetsExpected := []*budgetDom.Budget{
		{
			ID:        999,
			Name:      "Test 1",
			Year:      2023,
			Month:     5,
			ProjectId: baseProject.ID,
		},
	}
	suite.mockProjectRepo.On("FindById", ctx, baseProject.ID).Return(baseProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(1000), nil)
	suite.mockBudgetRepo.On("FindByProjectIds", ctx, []uint{baseProject.ID}).Return(budgetsExpected, nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", ctx, mock.Anything).Return(uint(0), gorm.ErrInvalidField)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestCloneErrorCreateLog() {
	require := require.New(suite.T())
	ctx := context.Background()
	baseProject := &domain.Project{
		ID:     suite.cloneId,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	budgetsExpected := []*budgetDom.Budget{
		{
			ID:        999,
			Name:      "Test 1",
			Year:      2023,
			Month:     5,
			ProjectId: baseProject.ID,
		},
	}
	suite.mockProjectRepo.On("FindById", ctx, baseProject.ID).Return(baseProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(uint(1000), nil)
	suite.mockBudgetRepo.On("FindByProjectIds", ctx, []uint{baseProject.ID}).Return(budgetsExpected, nil)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Create", ctx, mock.Anything).Return(uint(0), nil)
	suite.mockProjectLogRepo.On("WithTransaction", nil).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Create", ctx, mock.Anything).Return(uint(0), gorm.ErrInvalidField)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestFindByUserSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	projectsExpected := []*domain.Project{
		{
			ID:     999,
			UserId: suite.userId,
			Type:   suite.typeBudget,
		},
		{
			ID:     1000,
			UserId: suite.userId,
			Type:   suite.typeBudget,
		},
	}
	projectIds := []uint{999, 1000}
	budgetsExpected := []*budgetDom.Budget{
		{
			ID:        999,
			Name:      "Test 1",
			Year:      2023,
			Month:     5,
			ProjectId: projectsExpected[0].ID,
		},
		{
			ID:        1000,
			Name:      "Test 2",
			Year:      2023,
			Month:     5,
			ProjectId: projectsExpected[1].ID,
		},
	}
	suite.mockProjectRepo.On("FindByUserId", ctx, suite.userId).Return(projectsExpected, nil)
	suite.mockBudgetRepo.On("FindByProjectIds", ctx, projectIds).Return(budgetsExpected, nil)

	res, err := suite.app.FindByUser(ctx, suite.userId)

	require.NoError(err)
	require.Equal(len(projectsExpected), len(res))
	require.Equal(projectsExpected[0].ID, res[0].ID)
	require.Equal(projectsExpected[0].Type, res[0].Type)
	require.Equal(budgetsExpected[0].Name, res[0].Name)
	require.Equal(budgetsExpected[0].Year, res[0].Data["year"])
	require.Equal(budgetsExpected[0].Month, res[0].Data["month"])
	require.Equal(projectsExpected[1].ID, res[1].ID)
	require.Equal(projectsExpected[1].Type, res[1].Type)
	require.Equal(budgetsExpected[1].Name, res[1].Name)
	require.Equal(budgetsExpected[1].Year, res[1].Data["year"])
	require.Equal(budgetsExpected[1].Month, res[1].Data["month"])
}

func (suite *TestSuite) TestFindByUserErrorFindProjects() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockProjectRepo.On("FindByUserId", ctx, suite.userId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindByUser(ctx, suite.userId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestSuite) TestFindByUserErrorFindBudgets() {
	require := require.New(suite.T())
	ctx := context.Background()
	projectsExpected := []*domain.Project{
		{
			ID:     999,
			UserId: suite.userId,
			Type:   suite.typeBudget,
		},
		{
			ID:     1000,
			UserId: suite.userId,
			Type:   suite.typeBudget,
		},
	}
	projectIds := []uint{999, 1000}
	suite.mockProjectRepo.On("FindByUserId", ctx, suite.userId).Return(projectsExpected, nil)
	suite.mockBudgetRepo.On("FindByProjectIds", ctx, projectIds).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindByUser(ctx, suite.userId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestSuite) TestFindLogsByProjectSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	detail := `{"cloneId": 1}`
	logsExpected := []*domain.ProjectLog{
		{
			ID:          999,
			Description: "Test",
			ProjectId:   suite.cloneId,
		},
		{
			ID:          1000,
			Description: "Test",
			Detail:      &detail,
			ProjectId:   suite.cloneId,
		},
	}
	suite.mockProjectLogRepo.On("FindByProjectId", ctx, suite.cloneId).Return(logsExpected, nil)

	res, err := suite.app.FindLogsByProject(ctx, suite.cloneId)

	require.NoError(err)
	require.Equal(logsExpected, res)
}

func (suite *TestSuite) TestFindLogsByProjectError() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockProjectLogRepo.On("FindByProjectId", ctx, suite.cloneId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindLogsByProject(ctx, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	existsProject := &domain.Project{
		ID:   999,
		Type: suite.typeBudget,
	}
	suite.mockProjectRepo.On("FindById", ctx, suite.cloneId).Return(existsProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("DeleteByProjectId", ctx, existsProject.ID).Return(nil)
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Delete", ctx, suite.cloneId).Return(nil)

	err := suite.app.Delete(ctx, suite.cloneId)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteErrorFind() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockProjectRepo.On("FindById", ctx, suite.cloneId).Return(nil, gorm.ErrRecordNotFound)

	err := suite.app.Delete(ctx, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
}

func (suite *TestSuite) TestDeleteErrorDeleteProject() {
	require := require.New(suite.T())
	ctx := context.Background()
	existsProject := &domain.Project{
		ID:   999,
		Type: suite.typeBudget,
	}
	suite.mockProjectRepo.On("FindById", ctx, suite.cloneId).Return(existsProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("DeleteByProjectId", ctx, existsProject.ID).Return(gorm.ErrInvalidField)

	err := suite.app.Delete(ctx, suite.cloneId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

package application

import (
	"api-your-accounts/project/domain"
	"api-your-accounts/project/domain/mocks"
	"api-your-accounts/shared/domain/persistent"
	mocksShared "api-your-accounts/shared/domain/persistent/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	userId                 uint
	typeBudget             domain.ProjectType
	cloneId                uint
	mockTransactionManager *mocksShared.TransactionManager
	mockProjectRepo        *mocks.ProjectRepository
	mockProjectLogRepo     *mocks.ProjectLogRepository
	app                    IProjectApp
}

func (suite *TestSuite) SetupSuite() {
	suite.userId = 1
	suite.typeBudget = domain.Budget
	suite.cloneId = 1
}

func (suite *TestSuite) SetupTest() {
	suite.mockTransactionManager = mocksShared.NewTransactionManager(suite.T())
	suite.mockProjectRepo = mocks.NewProjectRepository(suite.T())
	suite.mockProjectLogRepo = mocks.NewProjectLogRepository(suite.T())
	suite.app = NewProjectApp(suite.mockTransactionManager, suite.mockProjectRepo, suite.mockProjectLogRepo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	project := &domain.Project{
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	projectExpected := &domain.Project{
		ID:     999,
		UserId: project.UserId,
		Type:   project.Type,
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(projectExpected, nil)
	suite.mockProjectLogRepo.On("WithTransaction", nil).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Create", ctx, mock.Anything).Return(nil, nil)

	res, err := suite.app.Create(ctx, project, &suite.cloneId)

	require.NoError(err)
	require.NotEmpty(res.ID)
	require.Equal(projectExpected.ID, res.ID)
	require.Equal(projectExpected.UserId, res.UserId)
	require.Equal(projectExpected.Type, res.Type)
}

func (suite *TestSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	ctx := context.Background()
	project := &domain.Project{
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(ctx, project, nil)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestCreateErrorCreateProject() {
	require := require.New(suite.T())
	ctx := context.Background()
	project := &domain.Project{
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(nil, errExpected)

	res, err := suite.app.Create(ctx, project, nil)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	ctx := context.Background()
	project := &domain.Project{
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	projectExpected := &domain.Project{
		ID:     999,
		UserId: project.UserId,
		Type:   project.Type,
	}
	errExpected := errors.New("Error in creation project")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockProjectRepo.On("WithTransaction", nil).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Create", ctx, mock.Anything).Return(projectExpected, nil)
	suite.mockProjectLogRepo.On("WithTransaction", nil).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Create", ctx, mock.Anything).Return(nil, errExpected)

	res, err := suite.app.Create(ctx, project, nil)

	require.EqualError(errExpected, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestCloneSuccess() {
	require := require.New(suite.T())
	ctx := context.Background()
	baseProject := &domain.Project{
		ID:     suite.cloneId,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}
	suite.mockProjectRepo.On("FindById", ctx, suite.cloneId).Return(baseProject, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(nil)

	_, err := suite.app.Clone(ctx, suite.cloneId)

	require.NoError(err)
}

func (suite *TestSuite) TestCloneError() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockProjectRepo.On("FindById", ctx, suite.cloneId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.Clone(ctx, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Nil(res)
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
	suite.mockProjectRepo.On("FindByUserId", ctx, suite.userId).Return(projectsExpected, nil)

	res, err := suite.app.FindByUser(ctx, suite.userId)

	require.NoError(err)
	require.Equal(projectsExpected, res)
}

func (suite *TestSuite) TestFindByUserError() {
	require := require.New(suite.T())
	ctx := context.Background()
	suite.mockProjectRepo.On("FindByUserId", ctx, suite.userId).Return(nil, gorm.ErrRecordNotFound)

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
	suite.mockProjectRepo.On("Delete", ctx, suite.cloneId).Return(nil)

	err := suite.app.Delete(ctx, suite.cloneId)

	require.NoError(err)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

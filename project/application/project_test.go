package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/project/domain"
	"your-accounts-api/project/domain/mocks"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"

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
	mockTransactionManager *mocks_shared.TransactionManager
	mockTx                 *mocks_shared.Transaction
	mockProjectRepo        *mocks.ProjectRepository
	mockProjectLogRepo     *mocks.ProjectLogRepository
	app                    IProjectApp
	ctx                    context.Context
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.userId = 1
	suite.typeBudget = domain.TypeBudget
	suite.cloneId = 1
	suite.ctx = context.Background()
}

func (suite *TestSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockTx = mocks_shared.NewTransaction(suite.T())
	suite.mockProjectRepo = mocks.NewProjectRepository(suite.T())
	suite.mockProjectLogRepo = mocks.NewProjectLogRepository(suite.T())
	suite.app = NewProjectApp(suite.mockTransactionManager, suite.mockProjectRepo, suite.mockProjectLogRepo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockProjectRepo.On("WithTransaction", suite.mockTx).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Save", suite.ctx, mock.Anything).Return(uint(999), nil)

	res, err := suite.app.Create(suite.ctx, suite.userId, suite.typeBudget, suite.mockTx)

	require.NoError(err)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project")
	suite.mockProjectRepo.On("WithTransaction", suite.mockTx).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, suite.typeBudget, suite.mockTx)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	suite.mockProjectRepo.On("WithTransaction", suite.mockTx).Return(suite.mockProjectRepo)
	suite.mockProjectRepo.On("Delete", suite.ctx, suite.cloneId).Return(nil)

	err := suite.app.Delete(suite.ctx, suite.cloneId, suite.mockTx)

	require.NoError(err)
}

func (suite *TestSuite) TestCreateLogSuccess() {
	require := require.New(suite.T())
	suite.mockProjectLogRepo.On("WithTransaction", suite.mockTx).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)

	err := suite.app.CreateLog(suite.ctx, "Save", suite.cloneId, nil, suite.mockTx)

	require.NoError(err)
}

func (suite *TestSuite) TestCreateLogError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project")
	suite.mockProjectLogRepo.On("WithTransaction", suite.mockTx).Return(suite.mockProjectLogRepo)
	suite.mockProjectLogRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	err := suite.app.CreateLog(suite.ctx, "Save", suite.cloneId, nil, suite.mockTx)

	require.EqualError(errExpected, err.Error())
}

func (suite *TestSuite) TestFindLogsByProjectSuccess() {
	require := require.New(suite.T())
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
	suite.mockProjectLogRepo.On("SearchAllByExample", suite.ctx, domain.ProjectLog{
		ProjectId: suite.cloneId,
	}).Return(logsExpected, nil)

	res, err := suite.app.FindLogsByProject(suite.ctx, suite.cloneId)

	require.NoError(err)
	require.Equal(logsExpected, res)
}

func (suite *TestSuite) TestFindLogsByProjectError() {
	require := require.New(suite.T())
	suite.mockProjectLogRepo.On("SearchAllByExample", suite.ctx, domain.ProjectLog{
		ProjectId: suite.cloneId,
	}).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindLogsByProject(suite.ctx, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

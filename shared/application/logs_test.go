package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/mocks"
	"your-accounts-api/shared/domain/persistent"
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
	code                   domain.CodeLog
	cloneId                uint
	mockTransactionManager *mocks_shared.TransactionManager
	mockTx                 *mocks_shared.Transaction
	mockLogRepo            *mocks.LogRepository
	app                    ILogApp
	ctx                    context.Context
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.userId = 1
	suite.code = domain.Budget
	suite.cloneId = 1
	suite.ctx = context.Background()
}

func (suite *TestSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockTx = mocks_shared.NewTransaction(suite.T())
	suite.mockLogRepo = mocks.NewLogRepository(suite.T())
	suite.app = NewLogApp(suite.mockTransactionManager, suite.mockLogRepo)
}

func (suite *TestSuite) TestCreateLogSuccess() {
	require := require.New(suite.T())
	suite.mockLogRepo.On("WithTransaction", suite.mockTx).Return(suite.mockLogRepo)
	suite.mockLogRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)

	err := suite.app.Create(suite.ctx, "Save", suite.code, suite.cloneId, nil, suite.mockTx)

	require.NoError(err)
}

func (suite *TestSuite) TestCreateLogError() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project")
	suite.mockLogRepo.On("WithTransaction", suite.mockTx).Return(suite.mockLogRepo)
	suite.mockLogRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	err := suite.app.Create(suite.ctx, "Save", suite.code, suite.cloneId, nil, suite.mockTx)

	require.EqualError(errExpected, err.Error())
}

func (suite *TestSuite) TestFindByProjectSuccess() {
	require := require.New(suite.T())
	detail := map[string]any{
		"cloneId": 1,
	}
	logsExpected := []domain.Log{
		{
			ID:          999,
			Description: "Test",
			ResourceId:  suite.cloneId,
		},
		{
			ID:          1000,
			Description: "Test",
			Detail:      detail,
			ResourceId:  suite.cloneId,
		},
	}
	suite.mockLogRepo.On("SearchAllByExample", suite.ctx, domain.Log{
		Code:       suite.code,
		ResourceId: suite.cloneId,
	}).Return(logsExpected, nil)

	res, err := suite.app.FindByProject(suite.ctx, suite.code, suite.cloneId)

	require.NoError(err)
	require.Equal(logsExpected, res)
}

func (suite *TestSuite) TestFindByProjectError() {
	require := require.New(suite.T())
	suite.mockLogRepo.On("SearchAllByExample", suite.ctx, domain.Log{
		Code:       suite.code,
		ResourceId: suite.cloneId,
	}).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindByProject(suite.ctx, suite.code, suite.cloneId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestSuite) TestDeleteOrphanSuccess() {
	require := require.New(suite.T())
	suite.mockLogRepo.On("DeleteByResourceIdNotExists", suite.ctx).Return(nil)

	err := suite.app.DeleteOrphan(suite.ctx)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteOrphanError() {
	require := require.New(suite.T())
	suite.mockLogRepo.On("DeleteByResourceIdNotExists", suite.ctx).Return(gorm.ErrRecordNotFound)

	err := suite.app.DeleteOrphan(suite.ctx)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
}

func (suite *TestSuite) TestDeleteOldSuccess() {
	require := require.New(suite.T())
	resourceIdsExpected := []uint{999, 1000}
	suite.mockLogRepo.On("SearchResourceIdsWithLimitExceeded", suite.ctx).Return(resourceIdsExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogRepo.On("WithTransaction", nil).Return(suite.mockLogRepo)
	suite.mockLogRepo.On("DeleteByResourceIdAndIdLessThanLimit", suite.ctx, mock.Anything).Return(nil).Times(2)

	err := suite.app.DeleteOld(suite.ctx)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteOldSuccessWithoutRecords() {
	require := require.New(suite.T())
	suite.mockLogRepo.On("SearchResourceIdsWithLimitExceeded", suite.ctx).Return([]uint{}, nil)

	err := suite.app.DeleteOld(suite.ctx)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteOldErrorSearch() {
	require := require.New(suite.T())
	suite.mockLogRepo.On("SearchResourceIdsWithLimitExceeded", suite.ctx).Return(nil, gorm.ErrRecordNotFound)

	err := suite.app.DeleteOld(suite.ctx)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
}

func (suite *TestSuite) TestDeleteOldErrorTransaction() {
	require := require.New(suite.T())
	resourceIdsExpected := []uint{999, 1000}
	errExpected := errors.New("Error in transaction")
	suite.mockLogRepo.On("SearchResourceIdsWithLimitExceeded", suite.ctx).Return(resourceIdsExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	err := suite.app.DeleteOld(suite.ctx)

	require.EqualError(errExpected, err.Error())
}

func (suite *TestSuite) TestDeleteOldErrorDelete() {
	require := require.New(suite.T())
	resourceIdsExpected := []uint{999, 1000}
	errExpected := errors.New("Error in transaction")
	suite.mockLogRepo.On("SearchResourceIdsWithLimitExceeded", suite.ctx).Return(resourceIdsExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockLogRepo.On("WithTransaction", nil).Return(suite.mockLogRepo)
	suite.mockLogRepo.On("DeleteByResourceIdAndIdLessThanLimit", suite.ctx, mock.Anything).Return(errExpected)

	err := suite.app.DeleteOld(suite.ctx)

	require.EqualError(errExpected, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

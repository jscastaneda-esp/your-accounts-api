package project

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"your-accounts-api/project/domain"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"
	"your-accounts-api/shared/domain/test_utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuite struct {
	suite.Suite
	userId     uint
	typeBudget domain.ProjectType
	mock       sqlmock.Sqlmock
	mockTX     *mocks_shared.Transaction
	repository domain.ProjectRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.userId = 1
	suite.typeBudget = domain.TypeBudget

	require := require.New(suite.T())

	var (
		db  *sql.DB
		err error
	)

	db, suite.mock, err = sqlmock.New()
	require.NoError(err)

	DB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(err)

	suite.repository = NewRepository(DB)
}

func (suite *TestSuite) SetupTest() {
	suite.mockTX = mocks_shared.NewTransaction(suite.T())
}

func (suite *TestSuite) TearDownTest() {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestWithTransactionSuccessNew() {
	require := require.New(suite.T())

	suite.mockTX.On("Get").Return(new(gorm.DB))

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.NotEqual(suite.repository, repo)
}

func (suite *TestSuite) TestWithTransactionSuccessExists() {
	require := require.New(suite.T())

	suite.mockTX.On("Get").Return(new(sql.DB))

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.Equal(suite.repository, repo)
}

func (suite *TestSuite) TestSaveSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `projects` (`created_at`,`user_id`,`type`) VALUES (?,?,?)")).
		WithArgs(test_utils.AnyTime{}, suite.userId, suite.typeBudget).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	project := domain.Project{
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}

	res, err := suite.repository.Save(context.Background(), project)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `projects` (`created_at`,`user_id`,`type`) VALUES (?,?,?)")).
		WithArgs(test_utils.AnyTime{}, suite.userId, suite.typeBudget).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	project := domain.Project{
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}

	res, err := suite.repository.Save(context.Background(), project)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `project_logs` WHERE `project_logs`.`project_id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `projects` WHERE `projects`.`id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.Delete(context.Background(), id)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteError() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `project_logs` WHERE `project_logs`.`project_id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `projects` WHERE `projects`.`id` = ?")).
		WithArgs(id).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

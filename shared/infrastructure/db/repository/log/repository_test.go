package log

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	"your-accounts-api/shared/domain"
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
	description string
	detail      string
	code        domain.CodeLog
	resourceId  uint
	mock        sqlmock.Sqlmock
	mockTX      *mocks_shared.Transaction
	repository  domain.LogRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.description = "Test"
	suite.detail = `{"test":"test"}`
	suite.code = domain.Budget
	suite.resourceId = 1

	require := require.New(suite.T())

	var (
		db  *sql.DB
		err error
	)

	db, suite.mock, err = sqlmock.New()
	require.NoError(err)
	suite.mock.MatchExpectationsInOrder(false)

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
		ExpectExec(regexp.QuoteMeta("INSERT INTO `logs` (`created_at`,`description`,`detail`,`code`,`resource_id`) VALUES (?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, suite.description, suite.detail, suite.code, suite.resourceId).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	log := domain.Log{
		Description: suite.description,
		Detail:      &suite.detail,
		Code:        suite.code,
		ResourceId:  suite.resourceId,
	}

	res, err := suite.repository.Save(context.Background(), log)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `logs` (`created_at`,`description`,`detail`,`code`,`resource_id`) VALUES (?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, suite.description, nil, suite.code, suite.resourceId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	log := domain.Log{
		Description: suite.description,
		Code:        suite.code,
		ResourceId:  suite.resourceId,
	}

	res, err := suite.repository.Save(context.Background(), log)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSearchAllByExampleSuccess() {
	require := require.New(suite.T())
	example := domain.Log{
		Code:       suite.code,
		ResourceId: suite.resourceId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `logs` WHERE `logs`.`code` = ? AND `logs`.`resource_id` = ? ORDER BY created_at desc LIMIT 20")).
		WithArgs(suite.code, suite.resourceId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "description", "detail", "code", "resource_id"}).
			AddRow(999, time.Now(), suite.description, suite.detail, suite.code, suite.resourceId).
			AddRow(1000, time.Now(), suite.description, nil, suite.code, suite.resourceId),
		)

	projects, err := suite.repository.SearchAllByExample(context.Background(), example)

	require.NoError(err)
	require.NotNil(projects)
	require.Len(projects, 2)
	require.Equal(uint(999), projects[0].ID)
	require.Equal(suite.description, projects[0].Description)
	require.Equal(suite.detail, *projects[0].Detail)
	require.Equal(suite.code, projects[0].Code)
	require.Equal(suite.resourceId, projects[0].ResourceId)
	require.Equal(uint(1000), projects[1].ID)
	require.Equal(suite.description, projects[1].Description)
	require.Nil(projects[1].Detail)
	require.Equal(suite.code, projects[1].Code)
	require.Equal(suite.resourceId, projects[1].ResourceId)
}

func (suite *TestSuite) TestSearchAllByExampleError() {
	require := require.New(suite.T())
	example := domain.Log{
		Code:       suite.code,
		ResourceId: suite.resourceId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `logs` WHERE `logs`.`code` = ? AND `logs`.`resource_id` = ? ORDER BY created_at desc LIMIT 20")).
		WithArgs(suite.code, suite.resourceId).
		WillReturnError(gorm.ErrRecordNotFound)

	projects, err := suite.repository.SearchAllByExample(context.Background(), example)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(projects)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

package project_log

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	"your-accounts-api/project/domain"
	mocksShared "your-accounts-api/shared/domain/persistent/mocks"
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
	projectId   uint
	mock        sqlmock.Sqlmock
	mockTX      *mocksShared.Transaction
	repository  domain.ProjectLogRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.description = "Test"
	suite.detail = `{"test":"test"}`
	suite.projectId = 1

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
	suite.mockTX = mocksShared.NewTransaction(suite.T())
}

func (suite *TestSuite) TearDownTest() {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestWithTransactionSuccessNew() {
	require := require.New(suite.T())

	suite.mockTX.On("Get").Return(&gorm.DB{})

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.NotEqual(suite.repository, repo)
}

func (suite *TestSuite) TestWithTransactionSuccessExists() {
	require := require.New(suite.T())

	suite.mockTX.On("Get").Return(&sql.DB{})

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.Equal(suite.repository, repo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `project_logs` (`created_at`,`description`,`detail`,`project_id`) VALUES (?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, suite.description, suite.detail, suite.projectId).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	projectLog := domain.ProjectLog{
		Description: suite.description,
		Detail:      &suite.detail,
		ProjectId:   suite.projectId,
	}

	res, err := suite.repository.Create(context.Background(), projectLog)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res.ID)
	require.Equal(projectLog.Description, res.Description)
	require.Equal(projectLog.Detail, res.Detail)
	require.Equal(projectLog.ProjectId, res.ProjectId)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `project_logs` (`created_at`,`description`,`detail`,`project_id`) VALUES (?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, suite.description, nil, suite.projectId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	projectLog := domain.ProjectLog{
		Description: suite.description,
		ProjectId:   suite.projectId,
	}

	res, err := suite.repository.Create(context.Background(), projectLog)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestFindByProjectIdSuccess() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_logs` WHERE `project_logs`.`project_id` = ? ORDER BY created_at desc LIMIT 20")).
		WithArgs(suite.projectId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "description", "detail", "project_id"}).
			AddRow(999, time.Now(), suite.description, suite.detail, suite.projectId).
			AddRow(1000, time.Now(), suite.description, nil, suite.projectId),
		)

	projects, err := suite.repository.FindByProjectId(context.Background(), suite.projectId)

	require.NoError(err)
	require.NotNil(projects)
	require.Len(projects, 2)
	require.Equal(uint(999), projects[0].ID)
	require.Equal(suite.description, projects[0].Description)
	require.Equal(suite.detail, *projects[0].Detail)
	require.Equal(suite.projectId, projects[0].ProjectId)
	require.Equal(uint(1000), projects[1].ID)
	require.Equal(suite.description, projects[1].Description)
	require.Nil(projects[1].Detail)
	require.Equal(suite.projectId, projects[1].ProjectId)
}

func (suite *TestSuite) TestFindByProjectIdError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_logs` WHERE `project_logs`.`project_id` = ? ORDER BY created_at desc LIMIT 20")).
		WithArgs(suite.projectId).
		WillReturnError(gorm.ErrRecordNotFound)

	projects, err := suite.repository.FindByProjectId(context.Background(), suite.projectId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(projects)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

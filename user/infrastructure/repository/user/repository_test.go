package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	mocksShared "your-accounts-api/shared/domain/persistent/mocks"
	"your-accounts-api/shared/domain/test_utils"
	"your-accounts-api/user/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuite struct {
	suite.Suite
	uid                string
	email              string
	mock               sqlmock.Sqlmock
	mockTX             *mocksShared.Transaction
	repository         domain.UserRepository
	repositoryInstance domain.UserRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.uid = "test"
	suite.email = "example@exaple.com"

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

	suite.repository = newRepository(DB)
	suite.repositoryInstance = DefaultRepository()
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

func (suite *TestSuite) TestFindByUIDAndEmailSuccess() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`uid` = ? AND `users`.`email` = ? ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(suite.uid, suite.email).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "uid", "email"}).
			AddRow(999, time.Now(), time.Now(), suite.uid, suite.email),
		)

	user, err := suite.repository.FindByUIDAndEmail(context.Background(), suite.uid, suite.email)

	require.NoError(err)
	require.NotNil(user)
	require.Equal(suite.uid, user.UID)
	require.Equal(suite.email, user.Email)
}

func (suite *TestSuite) TestFindByUIDAndEmailError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`uid` = ? AND `users`.`email` = ? ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(suite.uid, suite.email).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := suite.repository.FindByUIDAndEmail(context.Background(), suite.uid, suite.email)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Nil(user)
}

func (suite *TestSuite) TestExistsByUIDSuccessTrue() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM users WHERE uid = ?")).
		WithArgs(suite.uid).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := suite.repository.ExistsByUID(context.Background(), suite.uid)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsByUIDSuccessFalse() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM users WHERE uid = ?")).
		WithArgs(suite.uid).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := suite.repository.ExistsByUID(context.Background(), suite.uid)

	require.NoError(err)
	require.False(exists)
}

func (suite *TestSuite) TestExistsByUIDError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM users WHERE uid = ?")).
		WithArgs(suite.uid).
		WillReturnError(gorm.ErrInvalidField)

	exists, err := suite.repository.ExistsByUID(context.Background(), suite.uid)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestExistsByEmailSuccessTrue() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM users WHERE email = ?")).
		WithArgs(suite.email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := suite.repository.ExistsByEmail(context.Background(), suite.email)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsByEmailSuccessFalse() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM users WHERE email = ?")).
		WithArgs(suite.email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := suite.repository.ExistsByEmail(context.Background(), suite.email)

	require.NoError(err)
	require.False(exists)
}

func (suite *TestSuite) TestExistsByEmailError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT COUNT(1) FROM users WHERE email = ?")).
		WithArgs(suite.email).
		WillReturnError(gorm.ErrInvalidField)

	exists, err := suite.repository.ExistsByEmail(context.Background(), suite.email)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`uid`,`email`) VALUES (?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.uid, suite.email).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}

	res, err := suite.repository.Create(context.Background(), user)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`uid`,`email`) VALUES (?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.uid, suite.email).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	user := domain.User{
		UID:   suite.uid,
		Email: suite.email,
	}

	res, err := suite.repository.Create(context.Background(), user)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSingleton() {
	require := require.New(suite.T())

	repository := DefaultRepository()

	require.Equal(suite.repositoryInstance, repository)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

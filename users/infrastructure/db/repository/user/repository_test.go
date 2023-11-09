package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"
	"your-accounts-api/shared/domain/test_utils"
	"your-accounts-api/users/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuite struct {
	suite.Suite
	email      string
	mock       sqlmock.Sqlmock
	mockTX     *mocks_shared.Transaction
	repository domain.UserRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.email = "example@exaple.com"

	require := require.New(suite.T())

	var (
		db  *sql.DB
		err error
	)

	db, suite.mock, err = sqlmock.New()
	require.NoError(err)
	suite.mock.MatchExpectationsInOrder(false)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
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
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","email") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, suite.email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(999)))
	suite.mock.ExpectCommit()
	user := domain.User{
		Email: suite.email,
	}

	res, err := suite.repository.Save(context.Background(), user)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","email") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, suite.email).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	user := domain.User{
		Email: suite.email,
	}

	res, err := suite.repository.Save(context.Background(), user)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSearchByExampleSuccess() {
	require := require.New(suite.T())
	example := domain.User{
		Email: suite.email,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."email" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(suite.email).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "email"}).
			AddRow(999, time.Now(), suite.email),
		)

	user, err := suite.repository.SearchByExample(context.Background(), example)

	require.NoError(err)
	require.NotNil(user)
	require.Equal(suite.email, user.Email)
}

func (suite *TestSuite) TestSearchByExampleError() {
	require := require.New(suite.T())
	example := domain.User{
		Email: suite.email,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."email" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(suite.email).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := suite.repository.SearchByExample(context.Background(), example)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Nil(user)
}

func (suite *TestSuite) TestExistsByExampleSuccessTrue() {
	require := require.New(suite.T())
	example := domain.User{
		Email: suite.email,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE "users"."email" = $1`)).
		WithArgs(suite.email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := suite.repository.ExistsByExample(context.Background(), example)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsByExampleSuccessFalse() {
	require := require.New(suite.T())
	example := domain.User{
		Email: suite.email,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE "users"."email" = $1`)).
		WithArgs(suite.email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := suite.repository.ExistsByExample(context.Background(), example)

	require.NoError(err)
	require.False(exists)
}

func (suite *TestSuite) TestExistsByExampleError() {
	require := require.New(suite.T())
	example := domain.User{
		Email: suite.email,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE "users"."email" = $1`)).
		WithArgs(suite.email).
		WillReturnError(gorm.ErrInvalidField)

	exists, err := suite.repository.ExistsByExample(context.Background(), example)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.False(exists)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

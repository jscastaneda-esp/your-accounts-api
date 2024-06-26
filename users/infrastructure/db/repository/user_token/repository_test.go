package user_token

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	mocks_persistent "your-accounts-api/mocks/shared/domain/persistent"
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
	token      string
	userId     uint
	expiresAt  time.Time
	mock       sqlmock.Sqlmock
	mockTX     *mocks_persistent.MockTransaction
	repository domain.UserTokenRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.token = "<token>"
	suite.userId = 999
	suite.expiresAt = time.Now().Add(1 * time.Hour)

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

	suite.mockTX = mocks_persistent.NewMockTransaction(suite.T())
	suite.repository = NewRepository(DB)
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

	getMock := suite.mockTX.On("Get").Return(new(sql.DB))

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.Equal(suite.repository, repo)
	getMock.Unset()
}

func (suite *TestSuite) TestSaveSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_tokens" ("created_at","token","user_id","expires_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, suite.token, suite.userId, suite.expiresAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	userToken := domain.UserToken{
		Token:     suite.token,
		UserId:    suite.userId,
		ExpiresAt: suite.expiresAt,
	}

	res, err := suite.repository.Save(context.Background(), userToken)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user_tokens" ("created_at","token","user_id","expires_at") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, suite.token, suite.userId, suite.expiresAt).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	userToken := domain.UserToken{
		Token:     suite.token,
		UserId:    suite.userId,
		ExpiresAt: suite.expiresAt,
	}

	res, err := suite.repository.Save(context.Background(), userToken)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSearchByExampleSuccess() {
	require := require.New(suite.T())
	example := domain.UserToken{
		Token:  suite.token,
		UserId: suite.userId,
	}
	userTokenExpected := domain.UserToken{
		ID:        999,
		Token:     suite.token,
		UserId:    suite.userId,
		ExpiresAt: suite.expiresAt,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_tokens" WHERE "user_tokens"."token" = $1 AND "user_tokens"."user_id" = $2 ORDER BY "user_tokens"."id" LIMIT 1`)).
		WithArgs(suite.token, suite.userId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "token", "user_id", "created_at", "expires_at"}).
			AddRow(userTokenExpected.ID, userTokenExpected.Token, userTokenExpected.UserId, time.Now(), userTokenExpected.ExpiresAt),
		)

	res, err := suite.repository.SearchByExample(context.Background(), example)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(userTokenExpected, res)
}

func (suite *TestSuite) TestSearchByExampleError() {
	require := require.New(suite.T())
	example := domain.UserToken{
		Token:  suite.token,
		UserId: suite.userId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_tokens" WHERE "user_tokens"."token" = $1 AND "user_tokens"."user_id" = $2 ORDER BY "user_tokens"."id" LIMIT 1`)).
		WithArgs(suite.token, suite.userId).
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := suite.repository.SearchByExample(context.Background(), example)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestDeleteByExpiresAtGreaterThanNowSuccess() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_tokens" WHERE expires_at < NOW()`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.DeleteByExpiresAtGreaterThanNow(context.Background())

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteByExpiresAtGreaterThanNowError() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "user_tokens" WHERE expires_at < NOW()`)).
		WillReturnError(gorm.ErrRecordNotFound)
	suite.mock.ExpectRollback()

	err := suite.repository.DeleteByExpiresAtGreaterThanNow(context.Background())

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

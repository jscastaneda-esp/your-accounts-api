package user_token

import (
	mocksShared "api-your-accounts/shared/domain/persistent/mocks"
	"api-your-accounts/user/domain"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

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
	mockTX     *mocksShared.Transaction
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

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(err)

	suite.mockTX = mocksShared.NewTransaction(suite.T())
	suite.repository = NewRepository(DB)
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

	getMock := suite.mockTX.On("Get").Return(&sql.DB{})

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.Equal(suite.repository, repo)
	getMock.Unset()
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "user_tokens" ("token","user_id","refreshed_by","expires_at","refreshed_at") 
		VALUES ($1,$2,$3,$4,$5) 
		RETURNING "id","created_at"
		`)).
		WithArgs(suite.token, suite.userId, nil, suite.expiresAt, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
	suite.mock.ExpectCommit()
	userToken := domain.UserToken{
		Token:     suite.token,
		UserId:    suite.userId,
		ExpiresAt: suite.expiresAt,
	}

	res, err := suite.repository.Create(context.Background(), userToken)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res.ID)
	require.Equal(userToken.Token, res.Token)
	require.Equal(userToken.UserId, res.UserId)
	require.Equal(userToken.ExpiresAt, res.ExpiresAt)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "user_tokens" ("token","user_id","refreshed_by","expires_at","refreshed_at") 
		VALUES ($1,$2,$3,$4,$5) 
		RETURNING "id","created_at"
		`)).
		WithArgs(suite.token, suite.userId, nil, suite.expiresAt, nil).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	userToken := domain.UserToken{
		Token:     suite.token,
		UserId:    suite.userId,
		ExpiresAt: suite.expiresAt,
	}

	res, err := suite.repository.Create(context.Background(), userToken)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestFindByTokenAndUserIdSuccess() {
	require := require.New(suite.T())
	userTokenExpected := &domain.UserToken{
		ID:        999,
		Token:     suite.token,
		UserId:    suite.userId,
		CreatedAt: time.Now(),
		ExpiresAt: suite.expiresAt,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "user_tokens"
		WHERE "user_tokens"."token" = $1
		AND "user_tokens"."user_id" = $2
		ORDER BY "user_tokens"."id" LIMIT 1
		`)).
		WithArgs(suite.token, suite.userId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "token", "user_id", "refreshed_by", "created_at", "expires_at", "refreshed_at"}).
			AddRow(userTokenExpected.ID, userTokenExpected.Token, userTokenExpected.UserId, nil, userTokenExpected.CreatedAt, userTokenExpected.ExpiresAt, nil),
		)

	res, err := suite.repository.FindByTokenAndUserId(context.Background(), suite.token, suite.userId)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(userTokenExpected, res)
}

func (suite *TestSuite) TestFindByTokenAndUserIdError() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "user_tokens"
		WHERE "user_tokens"."token" = $1
		AND "user_tokens"."user_id" = $2
		ORDER BY "user_tokens"."id" LIMIT 1
		`)).
		WithArgs(suite.token, suite.userId).
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := suite.repository.FindByTokenAndUserId(context.Background(), suite.token, suite.userId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestUpdateSuccess() {
	require := require.New(suite.T())
	userToken := domain.UserToken{
		ID:        999,
		Token:     suite.token,
		UserId:    suite.userId,
		CreatedAt: time.Now(),
		ExpiresAt: suite.expiresAt,
	}
	refreshedBy := uint(1000)
	refreshedAt := time.Now()
	userTokenExpected := domain.UserToken{
		ID:          userToken.ID,
		Token:       userToken.Token,
		UserId:      userToken.UserId,
		RefreshedBy: &refreshedBy,
		CreatedAt:   userToken.CreatedAt,
		ExpiresAt:   userToken.ExpiresAt,
		RefreshedAt: &refreshedAt,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "user_tokens"
		WHERE "user_tokens"."id" = $1
		ORDER BY "user_tokens"."id" LIMIT 1
		`)).
		WithArgs(userToken.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "token", "user_id", "refreshed_by", "created_at", "expires_at", "refreshed_at"}).
			AddRow(userToken.ID, userToken.Token, userToken.UserId, nil, userToken.CreatedAt, userToken.ExpiresAt, nil),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`
		UPDATE "user_tokens"
		SET "created_at"=$1,"token"=$2,"user_id"=$3,"refreshed_by"=$4,"expires_at"=$5,"refreshed_at"=$6
		WHERE "id" = $7
		`)).
		WithArgs(userTokenExpected.CreatedAt, userTokenExpected.Token, userTokenExpected.UserId, *userTokenExpected.RefreshedBy, userTokenExpected.ExpiresAt, *userTokenExpected.RefreshedAt, userTokenExpected.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	res, err := suite.repository.Update(context.Background(), userTokenExpected)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(&userTokenExpected, res)
}

func (suite *TestSuite) TestUpdateErrorFind() {
	require := require.New(suite.T())
	userToken := domain.UserToken{
		ID:        999,
		Token:     suite.token,
		UserId:    suite.userId,
		CreatedAt: time.Now(),
		ExpiresAt: suite.expiresAt,
	}
	refreshedBy := uint(1000)
	refreshedAt := time.Now()
	userTokenExpected := domain.UserToken{
		ID:          userToken.ID,
		Token:       userToken.Token,
		UserId:      userToken.UserId,
		RefreshedBy: &refreshedBy,
		CreatedAt:   userToken.CreatedAt,
		ExpiresAt:   userToken.ExpiresAt,
		RefreshedAt: &refreshedAt,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "user_tokens"
		WHERE "user_tokens"."id" = $1
		ORDER BY "user_tokens"."id" LIMIT 1
		`)).
		WithArgs(userToken.ID).
		WillReturnError(gorm.ErrInvalidField)

	res, err := suite.repository.Update(context.Background(), userTokenExpected)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestUpdateErrorSave() {
	require := require.New(suite.T())
	userToken := domain.UserToken{
		ID:        999,
		Token:     suite.token,
		UserId:    suite.userId,
		CreatedAt: time.Now(),
		ExpiresAt: suite.expiresAt,
	}
	refreshedBy := uint(1000)
	refreshedAt := time.Now()
	userTokenExpected := domain.UserToken{
		ID:          userToken.ID,
		Token:       userToken.Token,
		UserId:      userToken.UserId,
		RefreshedBy: &refreshedBy,
		CreatedAt:   userToken.CreatedAt,
		ExpiresAt:   userToken.ExpiresAt,
		RefreshedAt: &refreshedAt,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "user_tokens"
		WHERE "user_tokens"."id" = $1
		ORDER BY "user_tokens"."id" LIMIT 1
		`)).
		WithArgs(userToken.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "token", "user_id", "refreshed_by", "created_at", "expires_at", "refreshed_at"}).
			AddRow(userToken.ID, userToken.Token, userToken.UserId, nil, userToken.CreatedAt, userToken.ExpiresAt, nil),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`
		UPDATE "user_tokens"
		SET "created_at"=$1,"token"=$2,"user_id"=$3,"refreshed_by"=$4,"expires_at"=$5,"refreshed_at"=$6
		WHERE "id" = $7
		`)).
		WithArgs(userTokenExpected.CreatedAt, userTokenExpected.Token, userTokenExpected.UserId, *userTokenExpected.RefreshedBy, userTokenExpected.ExpiresAt, *userTokenExpected.RefreshedAt, userTokenExpected.ID).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	res, err := suite.repository.Update(context.Background(), userTokenExpected)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

package user_token

import (
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

	suite.repository = NewRepository(DB)
}

func (suite *TestSuite) TearDownTest() {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "user_tokens" ("token","user_id","refreshed_id","expires_at","refreshed_at") 
		VALUES ($1,$2,$3,$4,$5) 
		RETURNING "id","created_at"
		`)).
		WithArgs(suite.token, suite.userId, nil, suite.expiresAt, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
	suite.mock.ExpectCommit()
	userToken := &domain.UserToken{
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
		INSERT INTO "user_tokens" ("token","user_id","refreshed_id","expires_at","refreshed_at") 
		VALUES ($1,$2,$3,$4,$5) 
		RETURNING "id","created_at"
		`)).
		WithArgs(suite.token, suite.userId, nil, suite.expiresAt, nil).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	userToken := &domain.UserToken{
		Token:     suite.token,
		UserId:    suite.userId,
		ExpiresAt: suite.expiresAt,
	}

	res, err := suite.repository.Create(context.Background(), userToken)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

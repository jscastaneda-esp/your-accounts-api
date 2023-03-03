package repository

import (
	sharedDom "api-your-accounts/shared/domain"
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
	uuid       string
	email      string
	mock       sqlmock.Sqlmock
	repository domain.UserRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.uuid = "test"
	suite.email = "example@exaple.com"

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

	suite.repository = NewGORMRepository(DB)
}

func (suite *TestSuite) TearDownTest() {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestFindByUUIDAndEmailSuccess() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "users"
		WHERE "users"."uuid" = $1
		AND "users"."email" = $2
		ORDER BY "users"."id" LIMIT 1
		`)).
		WithArgs(suite.uuid, suite.email).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "uuid", "email"}).
			AddRow(999, time.Now(), time.Now(), suite.uuid, suite.email),
		)

	user, err := suite.repository.FindByUUIDAndEmail(context.Background(), suite.uuid, suite.email)
	require.NoError(err)
	require.NotNil(user)
	require.Equal(suite.uuid, user.UUID)
	require.Equal(suite.email, user.Email)
}

func (suite *TestSuite) TestFindByUUIDAndEmailError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "users"
		WHERE "users"."uuid" = $1
		AND "users"."email" = $2
		ORDER BY "users"."id" LIMIT 1
		`)).
		WithArgs(suite.uuid, suite.email).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := suite.repository.FindByUUIDAndEmail(context.Background(), suite.uuid, suite.email)
	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Nil(user)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "users" ("updated_at","uuid","email") 
		VALUES ($1,$2,$3) 
		RETURNING "id","created_at"
		`)).
		WithArgs(sharedDom.AnyTime{}, suite.uuid, suite.email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
	suite.mock.ExpectCommit()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}

	res, err := suite.repository.Create(context.Background(), user)
	require.NoError(err)
	require.NotNil(res)
	require.Equal(user.UUID, res.UUID)
	require.Equal(user.Email, res.Email)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "users" ("updated_at","uuid","email") 
		VALUES ($1,$2,$3) 
		RETURNING "id","created_at"
		`)).
		WithArgs(sharedDom.AnyTime{}, suite.uuid, suite.email).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	user := &domain.User{
		UUID:  suite.uuid,
		Email: suite.email,
	}

	res, err := suite.repository.Create(context.Background(), user)
	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

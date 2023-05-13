package persistent

import (
	"api-your-accounts/shared/domain/persistent"
	"api-your-accounts/shared/domain/persistent/mocks"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	mockTx *mocks.Transaction
	db     *gorm.DB
	tm     persistent.TransactionManager
}

func (suite *TestSuite) SetupSuite() {
	require := require.New(suite.T())

	var (
		db  *sql.DB
		err error
	)

	db, suite.mock, err = sqlmock.New()
	require.NoError(err)

	suite.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(err)

	suite.tm = NewTransactionManager(suite.db)
}

func (suite *TestSuite) SetupTest() {
	suite.mockTx = mocks.NewTransaction(suite.T())
}

func (suite *TestSuite) TearDownTest() {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestSetSuccess() {
	require := require.New(suite.T())
	tx := &gormTransaction{}

	err := tx.Set(suite.db)

	require.NoError(err)
}

func (suite *TestSuite) TestSetError() {
	require := require.New(suite.T())
	tx := &gormTransaction{}

	err := tx.Set(&sql.DB{})

	require.EqualError(ErrInvalidTX, err.Error())
}

func (suite *TestSuite) TestGetSuccess() {
	require := require.New(suite.T())
	tx := &gormTransaction{}
	tx.Set(suite.db)

	db := tx.Get()

	require.NotNil(db)
	require.IsType(suite.db, db)
	require.Equal(suite.db, db)
}

func (suite *TestSuite) TestTransactionSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO test(value) VALUES ($1)")).
		WithArgs("Test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("UPDATE test SET value = $1 WHERE id = $2")).
		WithArgs("TestU", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.tm.Transaction(func(tx persistent.Transaction) error {
		db, err := tx.Get().(*gorm.DB).DB()
		require.NoError(err)

		_, err = db.Exec("INSERT INTO test(value) VALUES ($1)", "Test")
		require.NoError(err)

		_, err = db.Exec("UPDATE test SET value = $1 WHERE id = $2", "TestU", 1)
		require.NoError(err)

		return nil
	})

	require.NoError(err)
}

func (suite *TestSuite) TestDefaultWithTransactionSuccessNew() {
	require := require.New(suite.T())
	suite.mockTx.On("Get").Return(suite.db)

	result := DefaultWithTransaction(suite.mockTx, func(db *gorm.DB) string {
		return "New"
	}, "Default")

	require.NotNil(result)
	require.Equal("New", result)
}

func (suite *TestSuite) TestDefaultWithTransactionSuccessDefault() {
	require := require.New(suite.T())
	suite.mockTx.On("Get").Return("other")

	result := DefaultWithTransaction(suite.mockTx, func(db *gorm.DB) string {
		return "New"
	}, "Default")

	require.NotNil(result)
	require.Equal("Default", result)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

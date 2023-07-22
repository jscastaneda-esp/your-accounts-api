package budget_available_balance

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"your-accounts-api/budget/domain"
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
	name       string
	amount     float64
	budgetId   uint
	mock       sqlmock.Sqlmock
	mockTX     *mocks_shared.Transaction
	repository domain.BudgetAvailableBalanceRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.amount = 10.0
	suite.budgetId = 1

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

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`) VALUES (?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	available := domain.BudgetAvailableBalance{
		Name:     suite.name,
		Amount:   suite.amount,
		BudgetId: suite.budgetId,
	}

	res, err := suite.repository.Create(context.Background(), available)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`) VALUES (?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	available := domain.BudgetAvailableBalance{
		Name:     suite.name,
		Amount:   suite.amount,
		BudgetId: suite.budgetId,
	}

	res, err := suite.repository.Create(context.Background(), available)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

package budget

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
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
	year       uint16
	month      uint8
	projectId  uint
	mock       sqlmock.Sqlmock
	mockTX     *mocks_shared.Transaction
	repository domain.BudgetRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.year = 2023
	suite.month = 1
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
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budgets` (`created_at`,`updated_at`,`name`,`year`,`month`,`fixed_income`,`additional_income`,`total_pending_payment`,`total_available_balance`,`pending_bills`,`total_balance`,`total`,`estimated_balance`,`total_payment`,`project_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), float64(0), float64(0), float64(0), suite.projectId).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	budget := domain.Budget{
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
	}

	res, err := suite.repository.Create(context.Background(), budget)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budgets` (`created_at`,`updated_at`,`name`,`year`,`month`,`fixed_income`,`additional_income`,`total_pending_payment`,`total_available_balance`,`pending_bills`,`total_balance`,`total`,`estimated_balance`,`total_payment`,`project_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), float64(0), float64(0), float64(0), suite.projectId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	budget := domain.Budget{
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
	}

	res, err := suite.repository.Create(context.Background(), budget)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestFinByIdSuccess() {
	require := require.New(suite.T())
	budgetExpected := domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budgets` WHERE `budgets`.`id` = ? ORDER BY `budgets`.`id` LIMIT 1")).
		WithArgs(budgetExpected.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending_payment", "total_available_balance", "pending_bills", "total_balance", "total", "estimated_balance", "total_payment", "project_id"}).
			AddRow(budgetExpected.ID, budgetExpected.CreatedAt, budgetExpected.UpdatedAt, budgetExpected.Name, budgetExpected.Year, budgetExpected.Month, budgetExpected.FixedIncome, budgetExpected.AdditionalIncome, budgetExpected.TotalPendingPayment, budgetExpected.TotalAvailableBalance, budgetExpected.PendingBills, budgetExpected.TotalBalance, budgetExpected.Total, budgetExpected.EstimatedBalance, budgetExpected.TotalPayment, budgetExpected.ProjectId),
		)

	budget, err := suite.repository.FindById(context.Background(), budgetExpected.ID)

	require.NoError(err)
	require.NotNil(budget)
	require.Equal(budgetExpected.ID, budget.ID)
	require.Equal(budgetExpected.Name, budget.Name)
	require.Equal(budgetExpected.Year, budget.Year)
	require.Equal(budgetExpected.Month, budget.Month)
	require.Equal(budgetExpected.ProjectId, budget.ProjectId)
}

func (suite *TestSuite) TestFinByIdError() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budgets` WHERE `budgets`.`id` = ? ORDER BY `budgets`.`id` LIMIT 1")).
		WithArgs(999).
		WillReturnError(gorm.ErrInvalidField)

	budget, err := suite.repository.FindById(context.Background(), 999)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(budget)
}

func (suite *TestSuite) TestFindByProjectIdsSuccess() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budgets` WHERE project_id IN (?)")).
		WithArgs(suite.projectId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending_payment", "total_available_balance", "pending_bills", "total_balance", "total", "estimated_balance", "total_payment", "project_id"}).
			AddRow(999, time.Now(), time.Now(), suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), float64(0), float64(0), float64(0), suite.projectId),
		)

	budgets, err := suite.repository.FindByProjectIds(context.Background(), []uint{suite.projectId})

	require.NoError(err)
	require.NotEmpty(budgets)
	require.Len(budgets, 1)
	require.Equal(uint(999), budgets[0].ID)
	require.Equal(suite.name, budgets[0].Name)
	require.Equal(suite.year, budgets[0].Year)
	require.Equal(suite.month, budgets[0].Month)
	require.Equal(suite.projectId, budgets[0].ProjectId)
}

func (suite *TestSuite) TestFindByProjectIdsError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budgets` WHERE project_id IN (?)")).
		WithArgs(suite.projectId).
		WillReturnError(gorm.ErrRecordNotFound)

	projects, err := suite.repository.FindByProjectIds(context.Background(), []uint{suite.projectId})

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(projects)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budget_available_balances` WHERE `budget_available_balances`.`budget_id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budget_bills` WHERE `budget_bills`.`budget_id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budgets` WHERE `budgets`.`id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.Delete(context.Background(), id)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteErrorDelete() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budget_available_balances` WHERE `budget_available_balances`.`budget_id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budget_bills` WHERE `budget_bills`.`budget_id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budgets` WHERE `budgets`.`id` = ?")).
		WithArgs(id).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

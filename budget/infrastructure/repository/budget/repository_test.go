package budget

import (
	"api-your-accounts/budget/domain"
	mocksShared "api-your-accounts/shared/domain/persistent/mocks"
	"api-your-accounts/shared/domain/test_utils"
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
	name       string
	year       uint16
	month      uint8
	projectId  uint
	mock       sqlmock.Sqlmock
	mockTX     *mocksShared.Transaction
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

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
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
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "budgets" ("updated_at","name","year","month","fixed_income","additional_income","total_pending_payment","total_available_balance","pending_bills","total_balance","total","estimated_balance","total_payment","project_id")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING "id","created_at"
		`)).
		WithArgs(test_utils.AnyTime{}, suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), float64(0), float64(0), float64(0), suite.projectId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
	suite.mock.ExpectCommit()
	budget := &domain.Budget{
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
	}

	res, err := suite.repository.Create(context.Background(), budget)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res.ID)
	require.Equal(budget.Name, res.Name)
	require.Equal(budget.Year, res.Year)
	require.Equal(budget.Month, res.Month)
	require.Equal(budget.ProjectId, res.ProjectId)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "budgets" ("updated_at","name","year","month","fixed_income","additional_income","total_pending_payment","total_available_balance","pending_bills","total_balance","total","estimated_balance","total_payment","project_id")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING "id","created_at"
		`)).
		WithArgs(test_utils.AnyTime{}, suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), float64(0), float64(0), float64(0), suite.projectId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	budget := &domain.Budget{
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
	}

	res, err := suite.repository.Create(context.Background(), budget)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestFinByIdSuccess() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
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
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
		WithArgs(999).
		WillReturnError(gorm.ErrInvalidField)

	budget, err := suite.repository.FindById(context.Background(), 999)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(budget)
}

func (suite *TestSuite) TestFindByProjectIdSuccess() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."project_id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
		WithArgs(suite.projectId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending_payment", "total_available_balance", "pending_bills", "total_balance", "total", "estimated_balance", "total_payment", "project_id"}).
			AddRow(999, time.Now(), time.Now(), suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), float64(0), float64(0), float64(0), suite.projectId),
		)

	budget, err := suite.repository.FindByProjectId(context.Background(), suite.projectId)

	require.NoError(err)
	require.NotNil(budget)
	require.Equal(uint(999), budget.ID)
	require.Equal(suite.name, budget.Name)
	require.Equal(suite.year, budget.Year)
	require.Equal(suite.month, budget.Month)
	require.Equal(suite.projectId, budget.ProjectId)
}

func (suite *TestSuite) TestFindByProjectIdError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."project_id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
		WithArgs(suite.projectId).
		WillReturnError(gorm.ErrRecordNotFound)

	projects, err := suite.repository.FindByProjectId(context.Background(), suite.projectId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(projects)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
		WithArgs(budgetExpected.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending_payment", "total_available_balance", "pending_bills", "total_balance", "total", "estimated_balance", "total_payment", "project_id"}).
			AddRow(budgetExpected.ID, budgetExpected.CreatedAt, budgetExpected.UpdatedAt, budgetExpected.Name, budgetExpected.Year, budgetExpected.Month, budgetExpected.FixedIncome, budgetExpected.AdditionalIncome, budgetExpected.TotalPendingPayment, budgetExpected.TotalAvailableBalance, budgetExpected.PendingBills, budgetExpected.TotalBalance, budgetExpected.Total, budgetExpected.EstimatedBalance, budgetExpected.TotalPayment, budgetExpected.ProjectId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`
		DELETE FROM "budgets"
		WHERE "budgets"."id" = $1
		`)).
		WithArgs(budgetExpected.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.Delete(context.Background(), budgetExpected.ID)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteErrorFind() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
		WithArgs(budgetExpected.ID).
		WillReturnError(gorm.ErrInvalidField)

	err := suite.repository.Delete(context.Background(), budgetExpected.ID)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func (suite *TestSuite) TestDeleteErrorDelete() {
	require := require.New(suite.T())
	budgetExpected := &domain.Budget{
		ID:        999,
		Name:      suite.name,
		Year:      suite.year,
		Month:     suite.month,
		ProjectId: suite.projectId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "budgets"
		WHERE "budgets"."id" = $1
		ORDER BY "budgets"."id" LIMIT 1
		`)).
		WithArgs(budgetExpected.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending_payment", "total_available_balance", "pending_bills", "total_balance", "total", "estimated_balance", "total_payment", "project_id"}).
			AddRow(budgetExpected.ID, budgetExpected.CreatedAt, budgetExpected.UpdatedAt, budgetExpected.Name, budgetExpected.Year, budgetExpected.Month, budgetExpected.FixedIncome, budgetExpected.AdditionalIncome, budgetExpected.TotalPendingPayment, budgetExpected.TotalAvailableBalance, budgetExpected.PendingBills, budgetExpected.TotalBalance, budgetExpected.Total, budgetExpected.EstimatedBalance, budgetExpected.TotalPayment, budgetExpected.ProjectId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`
		DELETE FROM "budgets"
		WHERE "budgets"."id" = $1
		`)).
		WithArgs(budgetExpected.ID).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), budgetExpected.ID)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

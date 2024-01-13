package budget

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	"your-accounts-api/budgets/domain"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"
	"your-accounts-api/shared/domain/test_utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuite struct {
	suite.Suite
	id         uint
	name       string
	year       uint16
	month      uint8
	userId     uint
	mock       sqlmock.Sqlmock
	mockTX     *mocks_shared.Transaction
	repository domain.BudgetRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.id = 999
	suite.name = "Test"
	suite.year = 2023
	suite.month = 1
	suite.userId = 1

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

func (suite *TestSuite) TestSaveNewSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budgets" ("created_at","updated_at","name","year","month","fixed_income","additional_income","total_pending","total_available","pending_bills","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, suite.userId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(999)))
	suite.mock.ExpectCommit()
	budget := domain.Budget{
		Name:   &suite.name,
		Year:   &suite.year,
		Month:  &suite.month,
		UserId: &suite.userId,
	}

	res, err := suite.repository.Save(context.Background(), budget)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestSaveExistsSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	zeroFloat := 0.0
	zeroUInt := uint8(0)
	budgetExpected := domain.Budget{
		ID:     &suite.id,
		Name:   &suite.name,
		Year:   &suite.year,
		Month:  &suite.month,
		UserId: &suite.userId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budgets" WHERE "budgets"."id" = $1 ORDER BY "budgets"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending", "total_available", "pending_bills", "user_id"}).
			AddRow(budgetExpected.ID, now.Add(-1*time.Hour), now, budgetExpected.Name, budgetExpected.Year, budgetExpected.Month, budgetExpected.FixedIncome, budgetExpected.AdditionalIncome, budgetExpected.TotalPending, budgetExpected.TotalAvailable, budgetExpected.PendingBills, budgetExpected.UserId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`UPDATE "budgets" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"year"=$4,"month"=$5,"fixed_income"=$6,"additional_income"=$7,"total_pending"=$8,"total_available"=$9,"pending_bills"=$10,"user_id"=$11 WHERE "id" = $12`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.year, suite.month, zeroFloat, zeroFloat, zeroFloat, zeroFloat, zeroUInt, suite.userId, suite.id).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	budget := domain.Budget{
		ID:               &suite.id,
		Name:             &suite.name,
		Year:             &suite.year,
		Month:            &suite.month,
		FixedIncome:      &zeroFloat,
		AdditionalIncome: &zeroFloat,
		TotalAvailable:   &zeroFloat,
		TotalPending:     &zeroFloat,
		PendingBills:     &zeroUInt,
		UserId:           &suite.userId,
	}

	res, err := suite.repository.Save(context.Background(), budget)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res)
}

func (suite *TestSuite) TestSaveExistsErrorFind() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budgets" WHERE "budgets"."id" = $1 ORDER BY "budgets"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)
	pendingBills := uint8(2)
	budget := domain.Budget{
		ID:           &suite.id,
		Name:         &suite.name,
		Year:         &suite.year,
		Month:        &suite.month,
		PendingBills: &pendingBills,
		UserId:       &suite.userId,
	}

	res, err := suite.repository.Save(context.Background(), budget)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budgets" ("created_at","updated_at","name","year","month","fixed_income","additional_income","total_pending","total_available","pending_bills","user_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, suite.userId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	budget := domain.Budget{
		Name:   &suite.name,
		Year:   &suite.year,
		Month:  &suite.month,
		UserId: &suite.userId,
	}

	res, err := suite.repository.Save(context.Background(), budget)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSearchSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	ids := []uint{1, 2}
	names := []string{"Test 1", "Test 2"}
	categories := []domain.BudgetBillCategory{domain.Education, domain.Financial}
	budgetExpected := domain.Budget{
		ID:     &suite.id,
		Name:   &suite.name,
		Year:   &suite.year,
		Month:  &suite.month,
		UserId: &suite.userId,
		BudgetAvailables: []domain.BudgetAvailable{
			{
				ID:       &ids[0],
				Name:     &names[0],
				BudgetId: &suite.id,
			},
			{
				ID:       &ids[1],
				Name:     &names[1],
				BudgetId: &suite.id,
			},
		},
		BudgetBills: []domain.BudgetBill{
			{
				ID:          &ids[0],
				Description: &names[0],
				Category:    &categories[0],
				BudgetId:    &suite.id,
			},
			{
				ID:          &ids[1],
				Description: &names[1],
				Category:    &categories[1],
				BudgetId:    &suite.id,
			},
		},
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."budget_id" = $1 ORDER BY budget_availables.created_at ASC`)).
		WithArgs(budgetExpected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(budgetExpected.BudgetAvailables[0].ID, now, now, budgetExpected.BudgetAvailables[0].Name, budgetExpected.BudgetAvailables[0].Amount, budgetExpected.BudgetAvailables[0].BudgetId).
			AddRow(budgetExpected.BudgetAvailables[1].ID, now, now, budgetExpected.BudgetAvailables[1].Name, budgetExpected.BudgetAvailables[1].Amount, budgetExpected.BudgetAvailables[1].BudgetId))
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."budget_id" = $1 ORDER BY budget_bills.created_at ASC`)).
		WithArgs(budgetExpected.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "description", "amount", "payment", "due_date", "complete", "budget_id", "category"}).
			AddRow(budgetExpected.BudgetBills[0].ID, now, now, budgetExpected.BudgetBills[0].Description, budgetExpected.BudgetBills[0].Amount, budgetExpected.BudgetBills[0].Payment, budgetExpected.BudgetBills[0].DueDate, budgetExpected.BudgetBills[0].Complete, budgetExpected.BudgetAvailables[0].BudgetId, budgetExpected.BudgetBills[0].Category).
			AddRow(budgetExpected.BudgetBills[1].ID, now, now, budgetExpected.BudgetBills[1].Description, budgetExpected.BudgetBills[1].Amount, budgetExpected.BudgetBills[1].Payment, budgetExpected.BudgetBills[1].DueDate, budgetExpected.BudgetBills[1].Complete, budgetExpected.BudgetAvailables[1].BudgetId, budgetExpected.BudgetBills[1].Category))
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budgets" WHERE "budgets"."id" = $1 ORDER BY "budgets"."id" LIMIT 1`)).
		WithArgs(budgetExpected.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending", "total_available", "pending_bills", "user_id"}).
			AddRow(budgetExpected.ID, now, now, budgetExpected.Name, budgetExpected.Year, budgetExpected.Month, budgetExpected.FixedIncome, budgetExpected.AdditionalIncome, budgetExpected.TotalPending, budgetExpected.TotalAvailable, budgetExpected.PendingBills, budgetExpected.UserId),
		)

	budget, err := suite.repository.Search(context.Background(), *budgetExpected.ID)

	require.NoError(err)
	require.NotNil(budget)
	require.Equal(budgetExpected.ID, budget.ID)
	require.Equal(budgetExpected.Name, budget.Name)
	require.Equal(budgetExpected.Year, budget.Year)
	require.Equal(budgetExpected.Month, budget.Month)
	require.Equal(budgetExpected.UserId, budget.UserId)
	require.Len(budget.BudgetAvailables, len(budgetExpected.BudgetAvailables))
	require.Len(budget.BudgetBills, len(budgetExpected.BudgetBills))
	require.Equal(budgetExpected.BudgetAvailables[0].ID, budget.BudgetAvailables[0].ID)
	require.Equal(budgetExpected.BudgetAvailables[0].Name, budget.BudgetAvailables[0].Name)
	require.Equal(budgetExpected.BudgetAvailables[1].ID, budget.BudgetAvailables[1].ID)
	require.Equal(budgetExpected.BudgetAvailables[1].Name, budget.BudgetAvailables[1].Name)
	require.Equal(budgetExpected.BudgetBills[0].ID, budget.BudgetBills[0].ID)
	require.Equal(budgetExpected.BudgetBills[0].Description, budget.BudgetBills[0].Description)
	require.Equal(budgetExpected.BudgetBills[0].Category, budget.BudgetBills[0].Category)
	require.Equal(budgetExpected.BudgetBills[1].ID, budget.BudgetBills[1].ID)
	require.Equal(budgetExpected.BudgetBills[1].Description, budget.BudgetBills[1].Description)
	require.Equal(budgetExpected.BudgetBills[1].Category, budget.BudgetBills[1].Category)
}

func (suite *TestSuite) TestSearchError() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budgets" WHERE "budgets"."id" = $1 ORDER BY "budgets"."id" LIMIT 1`)).
		WithArgs(999).
		WillReturnError(gorm.ErrInvalidField)

	budget, err := suite.repository.Search(context.Background(), 999)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(budget)
}

func (suite *TestSuite) TestSearchAllByExampleSuccess() {
	require := require.New(suite.T())
	example := domain.Budget{
		UserId: &suite.userId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budgets" WHERE "budgets"."user_id" = $1`)).
		WithArgs(suite.userId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "year", "month", "fixed_income", "additional_income", "total_pending", "total_available", "pending_bills", "total_saving", "user_id"}).
			AddRow(999, time.Now(), time.Now(), suite.name, suite.year, suite.month, float64(0), float64(0), float64(0), float64(0), 0, float64(0), suite.userId).
			AddRow(1000, time.Now(), time.Now(), suite.name, suite.year, suite.month+1, float64(0), float64(0), float64(0), float64(0), 0, float64(0), suite.userId+1),
		)

	budgets, err := suite.repository.SearchAllByExample(context.Background(), example)

	require.NoError(err)
	require.NotEmpty(budgets)
	require.Len(budgets, 2)
	require.Equal(uint(999), *budgets[0].ID)
	require.Equal(suite.name, *budgets[0].Name)
	require.Equal(suite.year, *budgets[0].Year)
	require.Equal(suite.month, *budgets[0].Month)
	require.Equal(suite.userId, *budgets[0].UserId)
	require.Equal(uint(1000), *budgets[1].ID)
	require.Equal(suite.name, *budgets[1].Name)
	require.Equal(suite.year, *budgets[1].Year)
	require.Equal(suite.month+1, *budgets[1].Month)
	require.Equal(suite.userId+1, *budgets[1].UserId)
}

func (suite *TestSuite) TestSearchAllByExampleError() {
	require := require.New(suite.T())
	example := domain.Budget{
		UserId: &suite.userId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budgets" WHERE "budgets"."user_id" = $1`)).
		WithArgs(suite.userId).
		WillReturnError(gorm.ErrRecordNotFound)

	projects, err := suite.repository.SearchAllByExample(context.Background(), example)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(projects)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_availables" WHERE "budget_availables"."budget_id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_bills" WHERE "budget_bills"."budget_id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budgets" WHERE "budgets"."id" = $1`)).
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
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_availables" WHERE "budget_availables"."budget_id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_bills" WHERE "budget_bills"."budget_id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budgets" WHERE "budgets"."id" = $1`)).
		WithArgs(id).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

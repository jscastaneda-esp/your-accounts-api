package budget_bill

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	"your-accounts-api/budgets/domain"
	mocks_persistent "your-accounts-api/mocks/shared/domain/persistent"
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
	id          uint
	description string
	amount      float64
	budgetId    uint
	category    domain.BudgetBillCategory
	mock        sqlmock.Sqlmock
	mockTX      *mocks_persistent.MockTransaction
	repository  domain.BudgetBillRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.id = 1
	suite.description = "Test"
	suite.amount = 10.0
	suite.budgetId = 1
	suite.category = domain.House

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
	suite.mockTX = mocks_persistent.NewMockTransaction(suite.T())
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

func (suite *TestSuite) TestSearchSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	billExpected := domain.BudgetBill{
		ID:          &suite.id,
		Description: &suite.description,
		Amount:      &suite.amount,
		BudgetId:    &suite.budgetId,
		Category:    &suite.category,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "description", "amount", "payment", "due_date", "complete", "budget_id", "category"}).
			AddRow(billExpected.ID, now.Add(-1*time.Hour), now, billExpected.Description, billExpected.Amount, float64(0), uint8(0), false, billExpected.BudgetId, billExpected.Category),
		)

	res, err := suite.repository.Search(context.Background(), suite.id)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(*billExpected.ID, *res.ID)
	require.Equal(*billExpected.Description, *res.Description)
	require.Equal(*billExpected.Amount, *res.Amount)
	require.Equal(*billExpected.BudgetId, *res.BudgetId)
	require.Equal(*billExpected.Category, *res.Category)
}

func (suite *TestSuite) TestSearchError() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)

	res, err := suite.repository.Search(context.Background(), suite.id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSaveNewSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_bills" ("created_at","updated_at","description","amount","payment","due_date","complete","budget_id","category") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(suite.id)))
	suite.mock.ExpectCommit()
	bill := domain.BudgetBill{
		Description: &suite.description,
		Amount:      &suite.amount,
		BudgetId:    &suite.budgetId,
		Category:    &suite.category,
	}

	res, err := suite.repository.Save(context.Background(), bill)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(suite.id, res)
}

func (suite *TestSuite) TestSaveExistsSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	dueDate := uint8(1)
	complete := true
	billExpected := domain.BudgetBill{
		ID:          &suite.id,
		Description: &suite.description,
		Amount:      &suite.amount,
		BudgetId:    &suite.budgetId,
		Category:    &suite.category,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "description", "amount", "payment", "due_date", "complete", "budget_id", "category"}).
			AddRow(billExpected.ID, now.Add(-1*time.Hour), now, billExpected.Description, billExpected.Amount, float64(0), uint8(0), false, billExpected.BudgetId, billExpected.Category),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`UPDATE "budget_bills" SET "created_at"=$1,"updated_at"=$2,"description"=$3,"amount"=$4,"payment"=$5,"due_date"=$6,"complete"=$7,"budget_id"=$8,"category"=$9 WHERE "id" = $10`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, suite.amount, dueDate, complete, suite.budgetId, suite.category, suite.id).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	bill := domain.BudgetBill{
		ID:          &suite.id,
		Description: &suite.description,
		Amount:      &suite.amount,
		Payment:     &suite.amount,
		DueDate:     &dueDate,
		Complete:    &complete,
		BudgetId:    &suite.budgetId,
		Category:    &suite.category,
	}

	res, err := suite.repository.Save(context.Background(), bill)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(suite.id, res)
}

func (suite *TestSuite) TestSaveExistsErrorFind() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)
	bill := domain.BudgetBill{
		ID:          &suite.id,
		Description: &suite.description,
		Amount:      &suite.amount,
		BudgetId:    &suite.budgetId,
		Category:    &suite.category,
	}

	res, err := suite.repository.Save(context.Background(), bill)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_bills" ("created_at","updated_at","description","amount","payment","due_date","complete","budget_id","category") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	bill := domain.BudgetBill{
		Description: &suite.description,
		Amount:      &suite.amount,
		BudgetId:    &suite.budgetId,
		Category:    &suite.category,
	}

	res, err := suite.repository.Save(context.Background(), bill)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSaveAllNewSuccess() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_bills" ("created_at","updated_at","description","amount","payment","due_date","complete","budget_id","category") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9),($10,$11,$12,$13,$14,$15,$16,$17,$18) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$19,"description"="excluded"."description","amount"="excluded"."amount","payment"="excluded"."payment","due_date"="excluded"."due_date","complete"="excluded"."complete","budget_id"="excluded"."budget_id","category"="excluded"."category" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category, test_utils.AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	bills := []domain.BudgetBill{
		{
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
		{
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
	}

	err := suite.repository.SaveAll(context.Background(), bills)

	require.NoError(err)
}

func (suite *TestSuite) TestSaveAllNewAndUpdateSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "description", "amount", "payment", "due_date", "complete", "budget_id", "category"}).
			AddRow(suite.id, now.Add(-1*time.Hour), now, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_bills" ("created_at","updated_at","description","amount","payment","due_date","complete","budget_id","category","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10),($11,$12,$13,$14,$15,$16,$17,$18,$19,DEFAULT) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$20,"description"="excluded"."description","amount"="excluded"."amount","payment"="excluded"."payment","due_date"="excluded"."due_date","complete"="excluded"."complete","budget_id"="excluded"."budget_id","category"="excluded"."category" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category, suite.id, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category, test_utils.AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	bills := []domain.BudgetBill{
		{
			ID:          &suite.id,
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
		{
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
	}

	err := suite.repository.SaveAll(context.Background(), bills)

	require.NoError(err)
}

func (suite *TestSuite) TestSaveAllUpdateSuccess() {
	require := require.New(suite.T())
	ids := []uint{suite.id, suite.id + 1}
	now := time.Now()
	dueDate := uint8(1)
	complete := true
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(ids[0]).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "description", "amount", "payment", "due_date", "complete", "budget_id", "category"}).
			AddRow(ids[0], now.Add(-1*time.Hour), now, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category),
		)
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(ids[1]).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "description", "amount", "payment", "due_date", "complete", "budget_id", "category"}).
			AddRow(ids[1], now.Add(-1*time.Hour), now, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_bills" ("created_at","updated_at","description","amount","payment","due_date","complete","budget_id","category","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10),($11,$12,$13,$14,$15,$16,$17,$18,$19,$20) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$21,"description"="excluded"."description","amount"="excluded"."amount","payment"="excluded"."payment","due_date"="excluded"."due_date","complete"="excluded"."complete","budget_id"="excluded"."budget_id","category"="excluded"."category" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, suite.amount, uint8(0), complete, suite.budgetId, suite.category, ids[0], test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), dueDate, false, suite.budgetId, suite.category, ids[1], test_utils.AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	bills := []domain.BudgetBill{
		{
			ID:          &ids[0],
			Description: &suite.description,
			Amount:      &suite.amount,
			Payment:     &suite.amount,
			Complete:    &complete,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
		{
			ID:          &ids[1],
			Description: &suite.description,
			Amount:      &suite.amount,
			DueDate:     &dueDate,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
	}

	err := suite.repository.SaveAll(context.Background(), bills)

	require.NoError(err)
}

func (suite *TestSuite) TestSaveAllErrorFind() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_bills" WHERE "budget_bills"."id" = $1 ORDER BY "budget_bills"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)
	bills := []domain.BudgetBill{
		{
			ID:          &suite.id,
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
		{
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
	}

	err := suite.repository.SaveAll(context.Background(), bills)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func (suite *TestSuite) TestSaveAllError() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_bills" ("created_at","updated_at","description","amount","payment","due_date","complete","budget_id","category") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9),($10,$11,$12,$13,$14,$15,$16,$17,$18) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$19,"description"="excluded"."description","amount"="excluded"."amount","payment"="excluded"."payment","due_date"="excluded"."due_date","complete"="excluded"."complete","budget_id"="excluded"."budget_id","category"="excluded"."category" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.description, suite.amount, float64(0), uint8(0), false, suite.budgetId, suite.category, test_utils.AnyTime{}).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	bills := []domain.BudgetBill{
		{
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
		{
			Description: &suite.description,
			Amount:      &suite.amount,
			BudgetId:    &suite.budgetId,
			Category:    &suite.category,
		},
	}

	err := suite.repository.SaveAll(context.Background(), bills)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_bills" WHERE "budget_bills"."id" = $1`)).
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
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_bills" WHERE "budget_bills"."id" = $1`)).
		WithArgs(id).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

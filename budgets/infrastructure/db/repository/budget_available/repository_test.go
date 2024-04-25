package budget_available

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
	id         uint
	name       string
	amount     float64
	budgetId   uint
	mock       sqlmock.Sqlmock
	mockTX     *mocks_persistent.MockTransaction
	repository domain.BudgetAvailableRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.id = 1
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

func (suite *TestSuite) TestSaveNewSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_availables" ("created_at","updated_at","name","amount","budget_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(suite.id)))
	suite.mock.ExpectCommit()
	available := domain.BudgetAvailable{
		Name:     &suite.name,
		Amount:   &suite.amount,
		BudgetId: &suite.budgetId,
	}

	res, err := suite.repository.Save(context.Background(), available)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(suite.id, res)
}

func (suite *TestSuite) TestSaveExistsSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	budgetExpected := domain.BudgetAvailable{
		ID:       &suite.id,
		Name:     &suite.name,
		Amount:   &suite.amount,
		BudgetId: &suite.budgetId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."id" = $1 ORDER BY "budget_availables"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(budgetExpected.ID, now.Add(-1*time.Hour), now, budgetExpected.Name, budgetExpected.Amount, budgetExpected.BudgetId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`UPDATE "budget_availables" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"amount"=$4,"budget_id"=$5 WHERE "id" = $6`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, suite.id).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	available := domain.BudgetAvailable{
		ID:       &suite.id,
		Name:     &suite.name,
		Amount:   &suite.amount,
		BudgetId: &suite.budgetId,
	}

	res, err := suite.repository.Save(context.Background(), available)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(suite.id, res)
}

func (suite *TestSuite) TestSaveExistsErrorFind() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."id" = $1 ORDER BY "budget_availables"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)
	available := domain.BudgetAvailable{
		ID:       &suite.id,
		Name:     &suite.name,
		Amount:   &suite.amount,
		BudgetId: &suite.budgetId,
	}

	res, err := suite.repository.Save(context.Background(), available)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSaveError() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_availables" ("created_at","updated_at","name","amount","budget_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	available := domain.BudgetAvailable{
		Name:     &suite.name,
		Amount:   &suite.amount,
		BudgetId: &suite.budgetId,
	}

	res, err := suite.repository.Save(context.Background(), available)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Zero(res)
}

func (suite *TestSuite) TestSaveAllNewSuccess() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_availables" ("created_at","updated_at","name","amount","budget_id") VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$11,"name"="excluded"."name","amount"="excluded"."amount","budget_id"="excluded"."budget_id" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	availables := []domain.BudgetAvailable{
		{
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
		{
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
	}

	err := suite.repository.SaveAll(context.Background(), availables)

	require.NoError(err)
}

func (suite *TestSuite) TestSaveAllNewAndUpdateSuccess() {
	require := require.New(suite.T())
	now := time.Now()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."id" = $1 ORDER BY "budget_availables"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(suite.id, now.Add(-1*time.Hour), now, suite.name, suite.amount, suite.budgetId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_availables" ("created_at","updated_at","name","amount","budget_id","id") VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,DEFAULT) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$12,"name"="excluded"."name","amount"="excluded"."amount","budget_id"="excluded"."budget_id" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, suite.id, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	availables := []domain.BudgetAvailable{
		{
			ID:       &suite.id,
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
		{
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
	}

	err := suite.repository.SaveAll(context.Background(), availables)

	require.NoError(err)
}

func (suite *TestSuite) TestSaveAllUpdateSuccess() {
	require := require.New(suite.T())
	ids := []uint{suite.id, suite.id + 1}
	now := time.Now()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."id" = $1 ORDER BY "budget_availables"."id" LIMIT 1`)).
		WithArgs(ids[0]).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(ids[0], now.Add(-1*time.Hour), now, suite.name, suite.amount, suite.budgetId),
		)
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."id" = $1 ORDER BY "budget_availables"."id" LIMIT 1`)).
		WithArgs(ids[1]).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(ids[1], now.Add(-1*time.Hour), now, suite.name, suite.amount, suite.budgetId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_availables" ("created_at","updated_at","name","amount","budget_id","id") VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$13,"name"="excluded"."name","amount"="excluded"."amount","budget_id"="excluded"."budget_id" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, ids[0], test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, ids[1], test_utils.AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(999)))
	suite.mock.ExpectCommit()
	availables := []domain.BudgetAvailable{
		{
			ID:       &ids[0],
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
		{
			ID:       &ids[1],
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
	}

	err := suite.repository.SaveAll(context.Background(), availables)

	require.NoError(err)
}

func (suite *TestSuite) TestSaveAllErrorFind() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "budget_availables" WHERE "budget_availables"."id" = $1 ORDER BY "budget_availables"."id" LIMIT 1`)).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)
	availables := []domain.BudgetAvailable{
		{
			ID:       &suite.id,
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
		{
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
	}

	err := suite.repository.SaveAll(context.Background(), availables)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func (suite *TestSuite) TestSaveAllError() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "budget_availables" ("created_at","updated_at","name","amount","budget_id") VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$11,"name"="excluded"."name","amount"="excluded"."amount","budget_id"="excluded"."budget_id" RETURNING "id"`)).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	availables := []domain.BudgetAvailable{
		{
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
		{
			Name:     &suite.name,
			Amount:   &suite.amount,
			BudgetId: &suite.budgetId,
		},
	}

	err := suite.repository.SaveAll(context.Background(), availables)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_availables" WHERE "budget_availables"."id" = $1`)).
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
		ExpectExec(regexp.QuoteMeta(`DELETE FROM "budget_availables" WHERE "budget_availables"."id" = $1`)).
		WithArgs(id).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

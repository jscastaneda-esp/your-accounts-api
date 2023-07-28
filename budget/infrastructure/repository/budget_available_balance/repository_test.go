package budget_available_balance

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
	id         uint
	name       string
	amount     float64
	budgetId   uint
	mock       sqlmock.Sqlmock
	mockTX     *mocks_shared.Transaction
	repository domain.BudgetAvailableBalanceRepository
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

func (suite *TestSuite) TestSaveNewSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`) VALUES (?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId).
		WillReturnResult(sqlmock.NewResult(int64(suite.id), 1))
	suite.mock.ExpectCommit()
	available := domain.BudgetAvailableBalance{
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
	budgetExpected := domain.BudgetAvailableBalance{
		ID:       &suite.id,
		Name:     &suite.name,
		Amount:   &suite.amount,
		BudgetId: &suite.budgetId,
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budget_available_balances` WHERE `budget_available_balances`.`id` = ? ORDER BY `budget_available_balances`.`id` LIMIT 1")).
		WithArgs(suite.id).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(budgetExpected.ID, now.Add(-1*time.Hour), now, budgetExpected.Name, budgetExpected.Amount, budgetExpected.BudgetId),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("UPDATE `budget_available_balances` SET `created_at`=?,`updated_at`=?,`name`=?,`amount`=?,`budget_id`=? WHERE `id` = ?")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, suite.id).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	available := domain.BudgetAvailableBalance{
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
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budget_available_balances` WHERE `budget_available_balances`.`id` = ? ORDER BY `budget_available_balances`.`id` LIMIT 1")).
		WithArgs(suite.id).
		WillReturnError(gorm.ErrInvalidField)
	available := domain.BudgetAvailableBalance{
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
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`) VALUES (?,?,?,?,?)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	available := domain.BudgetAvailableBalance{
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
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`) VALUES (?,?,?,?,?),(?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`name`=VALUES(`name`),`amount`=VALUES(`amount`),`budget_id`=VALUES(`budget_id`)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	availables := []domain.BudgetAvailableBalance{
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
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`,`id`) VALUES (?,?,?,?,?,?),(?,?,?,?,?,DEFAULT) ON DUPLICATE KEY UPDATE `updated_at`=?,`name`=VALUES(`name`),`amount`=VALUES(`amount`),`budget_id`=VALUES(`budget_id`)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, suite.id, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	availables := []domain.BudgetAvailableBalance{
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
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`,`id`) VALUES (?,?,?,?,?,?),(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`name`=VALUES(`name`),`amount`=VALUES(`amount`),`budget_id`=VALUES(`budget_id`)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, ids[0], test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, ids[1], test_utils.AnyTime{}).
		WillReturnResult(sqlmock.NewResult(int64(999), 1))
	suite.mock.ExpectCommit()
	availables := []domain.BudgetAvailableBalance{
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

func (suite *TestSuite) TestSaveAllError() {
	require := require.New(suite.T())
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("INSERT INTO `budget_available_balances` (`created_at`,`updated_at`,`name`,`amount`,`budget_id`) VALUES (?,?,?,?,?),(?,?,?,?,?) ON DUPLICATE KEY UPDATE `updated_at`=?,`name`=VALUES(`name`),`amount`=VALUES(`amount`),`budget_id`=VALUES(`budget_id`)")).
		WithArgs(test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}, test_utils.AnyTime{}, suite.name, suite.amount, suite.budgetId, test_utils.AnyTime{}).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	availables := []domain.BudgetAvailableBalance{
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

func (suite *TestSuite) TestSearchAllByExampleSuccess() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budget_available_balances` WHERE `budget_available_balances`.`budget_id` = ?")).
		WithArgs(suite.budgetId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "amount", "budget_id"}).
			AddRow(999, time.Now(), time.Now(), suite.name, suite.amount, suite.budgetId).
			AddRow(1000, time.Now(), time.Now(), suite.name, suite.amount+1, suite.budgetId+1),
		)
	example := domain.BudgetAvailableBalance{
		BudgetId: &suite.budgetId,
	}

	availables, err := suite.repository.SearchAllByExample(context.Background(), example)

	require.NoError(err)
	require.NotEmpty(availables)
	require.Len(availables, 2)
	require.Equal(uint(999), *availables[0].ID)
	require.Equal(suite.name, *availables[0].Name)
	require.Equal(suite.amount, *availables[0].Amount)
	require.Equal(suite.budgetId, *availables[0].BudgetId)
	require.Equal(uint(1000), *availables[1].ID)
	require.Equal(suite.name, *availables[1].Name)
	require.Equal(suite.amount+1, *availables[1].Amount)
	require.Equal(suite.budgetId+1, *availables[1].BudgetId)
}

func (suite *TestSuite) TestSearchAllByExampleError() {
	require := require.New(suite.T())
	suite.mock.
		ExpectQuery(regexp.QuoteMeta("SELECT * FROM `budget_available_balances` WHERE `budget_available_balances`.`budget_id` = ?")).
		WithArgs(suite.budgetId).
		WillReturnError(gorm.ErrInvalidField)
	example := domain.BudgetAvailableBalance{
		BudgetId: &suite.budgetId,
	}

	availables, err := suite.repository.SearchAllByExample(context.Background(), example)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Empty(availables)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	id := uint(999)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budget_available_balances` WHERE `budget_available_balances`.`id` = ?")).
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
		ExpectExec(regexp.QuoteMeta("DELETE FROM `budget_available_balances` WHERE `budget_available_balances`.`id` = ?")).
		WithArgs(id).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), id)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

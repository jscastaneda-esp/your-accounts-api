package application

import (
	"context"
	"errors"
	"testing"
	"your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/domain/mocks"
	mocks_logs "your-accounts-api/shared/application/mocks"
	shared "your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
	mocks_shared "your-accounts-api/shared/domain/persistent/mocks"
	"your-accounts-api/shared/domain/utils/convert"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TestBudgetSuite struct {
	suite.Suite
	budgetId                uint
	userId                  uint
	mockTransactionManager  *mocks_shared.TransactionManager
	mockBudgetRepo          *mocks.BudgetRepository
	mockBudgetAvailableRepo *mocks.BudgetAvailableRepository
	mockBudgetBillRepo      *mocks.BudgetBillRepository
	mockLogApp              *mocks_logs.ILogApp
	app                     IBudgetApp
	ctx                     context.Context
}

func (suite *TestBudgetSuite) SetupSuite() {
	suite.budgetId = 1
	suite.userId = 2
	suite.ctx = context.Background()
}

func (suite *TestBudgetSuite) SetupTest() {
	suite.mockTransactionManager = mocks_shared.NewTransactionManager(suite.T())
	suite.mockBudgetRepo = mocks.NewBudgetRepository(suite.T())
	suite.mockBudgetAvailableRepo = mocks.NewBudgetAvailableRepository(suite.T())
	suite.mockBudgetBillRepo = mocks.NewBudgetBillRepository(suite.T())
	suite.mockLogApp = mocks_logs.NewILogApp(suite.T())
	suite.app = NewBudgetApp(suite.mockTransactionManager, suite.mockBudgetRepo, suite.mockBudgetAvailableRepo, suite.mockBudgetBillRepo, suite.mockLogApp)
}

func (suite *TestBudgetSuite) TestCreateSuccess() {
	require := require.New(suite.T())
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetSuite) TestCreateErrorCreateLog() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation project log")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCreateErrorSave() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in creation budget")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCreateErrorTransaction() {
	require := require.New(suite.T())
	errExpected := errors.New("Error in transaction")
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Create(suite.ctx, suite.userId, "Test")

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneSuccess() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	category := domain.Education
	budgetExpected := &domain.Budget{
		ID:     &baseId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
		BudgetAvailables: []domain.BudgetAvailable{
			{
				ID:       &baseId,
				Name:     &name,
				BudgetId: &baseId,
			},
		},
		BudgetBills: []domain.BudgetBill{
			{
				ID:          &baseId,
				Description: &name,
				Category:    &category,
				BudgetId:    &baseId,
			},
		},
	}
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo)
	suite.mockBudgetAvailableRepo.On("SaveAll", suite.ctx, mock.Anything).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("SaveAll", suite.ctx, mock.Anything).Return(nil)
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.NoError(err)
	require.Equal(suite.budgetId, res)
}

func (suite *TestBudgetSuite) TestCloneErrorSearch() {
	require := require.New(suite.T())
	baseId := uint(999)
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(nil, errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorSaveAvailables() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &baseId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	errExpected := errors.New("Error in creation availables")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo)
	suite.mockBudgetAvailableRepo.On("SaveAll", suite.ctx, mock.Anything).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorSaveBills() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &baseId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	errExpected := errors.New("Error in creation bills")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo)
	suite.mockBudgetAvailableRepo.On("SaveAll", suite.ctx, mock.Anything).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("SaveAll", suite.ctx, mock.Anything).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorCreateLog() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &baseId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	errExpected := errors.New("Error in creation project log")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(suite.budgetId, nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo)
	suite.mockBudgetAvailableRepo.On("SaveAll", suite.ctx, mock.Anything).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo)
	suite.mockBudgetBillRepo.On("SaveAll", suite.ctx, mock.Anything).Return(nil)
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorSave() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &baseId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	errExpected := errors.New("Error in creation budget")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestCloneErrorTransaction() {
	require := require.New(suite.T())
	baseId := uint(999)
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &baseId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	errExpected := errors.New("Error in transaction")
	suite.mockBudgetRepo.On("Search", suite.ctx, baseId).Return(budgetExpected, nil)
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(errExpected)

	res, err := suite.app.Clone(suite.ctx, suite.userId, baseId)

	require.EqualError(errExpected, err.Error())
	require.Zero(res)
}

func (suite *TestBudgetSuite) TestFindByIdSuccess() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	ids := []uint{1, 2}
	names := []string{"Test 1", "Test 2"}
	categories := []domain.BudgetBillCategory{domain.Education, domain.Financial}
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
		BudgetAvailables: []domain.BudgetAvailable{
			{
				ID:       &ids[0],
				Name:     &names[0],
				BudgetId: &suite.budgetId,
			},
			{
				ID:       &ids[1],
				Name:     &names[1],
				BudgetId: &suite.budgetId,
			},
		},
		BudgetBills: []domain.BudgetBill{
			{
				ID:          &ids[0],
				Description: &names[0],
				Category:    &categories[0],
				BudgetId:    &suite.budgetId,
			},
			{
				ID:          &ids[1],
				Description: &names[1],
				Category:    &categories[1],
				BudgetId:    &suite.budgetId,
			},
		},
	}
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)

	res, err := suite.app.FindById(suite.ctx, suite.budgetId)

	require.NoError(err)
	require.Equal(budgetExpected, res)
}

func (suite *TestBudgetSuite) TestFindByIdError() {
	require := require.New(suite.T())
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(nil, gorm.ErrRecordNotFound)

	res, err := suite.app.FindById(suite.ctx, suite.budgetId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(res)
}

func (suite *TestBudgetSuite) TestFindByUserIdSuccess() {
	require := require.New(suite.T())
	example := domain.Budget{
		UserId: &suite.userId,
	}
	ids := []uint{uint(999), uint(1000)}
	names := []string{"Test 1", "Test 2"}
	year := uint16(2003)
	month := uint8(5)
	userIds := []uint{999, 1000}
	budgetsExpected := []domain.Budget{
		{
			ID:     &ids[0],
			Name:   &names[0],
			Year:   &year,
			Month:  &month,
			UserId: &userIds[0],
		},
		{
			ID:     &ids[1],
			Name:   &names[1],
			Year:   &year,
			Month:  &month,
			UserId: &userIds[1],
		},
	}
	suite.mockBudgetRepo.On("SearchAllByExample", suite.ctx, example).Return(budgetsExpected, nil)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.NoError(err)
	require.Equal(budgetsExpected, res)
}

func (suite *TestBudgetSuite) TestFindByUserIdError() {
	require := require.New(suite.T())
	example := domain.Budget{
		UserId: &suite.userId,
	}
	suite.mockBudgetRepo.On("SearchAllByExample", suite.ctx, example).Return(nil, gorm.ErrInvalidField)

	res, err := suite.app.FindByUserId(suite.ctx, suite.userId)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestBudgetSuite) TestChangesSuccess() {
	require := require.New(suite.T())
	changes := []Change{
		{
			ID:      uint(1),
			Section: domain.Main,
			Action:  shared.Update,
			Detail: map[string]any{
				"name":             "Test 1",
				"year":             2023,
				"month":            1,
				"fixedIncome":      10000,
				"additionalIncome": 0,
			},
		},
		{
			ID:      uint(1),
			Section: domain.Available,
			Action:  shared.Update,
			Detail: map[string]any{
				"name":   "Test 1",
				"amount": 10000,
			},
		},
		{
			ID:      uint(2),
			Section: domain.Available,
			Action:  shared.Delete,
			Detail: map[string]any{
				"name": "Test 2",
			},
		},
		{
			ID:      uint(1),
			Section: domain.Bill,
			Action:  shared.Update,
			Detail: map[string]any{
				"description": "Test 1",
				"amount":      8000,
				"dueDate":     1,
				"complete":    false,
				"category":    "education",
			},
		},
		{
			ID:      uint(2),
			Section: domain.Bill,
			Action:  shared.Delete,
			Detail: map[string]any{
				"description": "Test 2",
			},
		},
	}
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	ids := []uint{1, 2}
	names := []string{"Test 1", "Test 2"}
	categories := []domain.BudgetBillCategory{domain.Education, domain.Financial}
	zeroFloat := 0.0
	zeroBool := false
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
		BudgetAvailables: []domain.BudgetAvailable{
			{
				ID:       &ids[0],
				Name:     &names[0],
				Amount:   &zeroFloat,
				BudgetId: &suite.budgetId,
			},
			{
				ID:       &ids[1],
				Name:     &names[1],
				Amount:   &zeroFloat,
				BudgetId: &suite.budgetId,
			},
		},
		BudgetBills: []domain.BudgetBill{
			{
				ID:          &ids[0],
				Description: &names[0],
				Amount:      &zeroFloat,
				Payment:     &zeroFloat,
				Complete:    &zeroBool,
				Category:    &categories[0],
				BudgetId:    &suite.budgetId,
			},
			{
				ID:          &ids[1],
				Description: &names[1],
				Amount:      &zeroFloat,
				Payment:     &zeroFloat,
				Complete:    &zeroBool,
				Category:    &categories[1],
				BudgetId:    &suite.budgetId,
			},
		},
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	}).Times(5)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo).Times(2)
	suite.mockBudgetAvailableRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)
	suite.mockBudgetAvailableRepo.On("Delete", suite.ctx, mock.Anything).Return(nil)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo).Times(2)
	suite.mockBudgetBillRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), nil)
	suite.mockBudgetBillRepo.On("Delete", suite.ctx, mock.Anything).Return(nil)
	suite.mockLogApp.On("Create", suite.ctx, mock.Anything, shared.Budget, suite.budgetId, mock.Anything, nil).Return(nil).Times(5)
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)

	results := suite.app.Changes(suite.ctx, suite.budgetId, changes)

	require.NotEmpty(results)
	require.Len(results, 5)

	for _, result := range results {
		require.NoError(result.Err)
	}
}

func (suite *TestBudgetSuite) TestChangesErrorData() {
	require := require.New(suite.T())
	changes := []Change{
		{
			ID:      uint(1),
			Section: domain.Main,
			Action:  shared.Update,
		},
		{
			ID:      uint(1),
			Section: domain.Available,
			Action:  shared.Update,
		},
		{
			ID:      uint(2),
			Section: domain.Available,
			Action:  shared.Delete,
		},
		{
			ID:      uint(1),
			Section: domain.Bill,
			Action:  shared.Update,
		},
		{
			ID:      uint(2),
			Section: domain.Bill,
			Action:  shared.Delete,
		},
	}

	results := suite.app.Changes(suite.ctx, suite.budgetId, changes)

	require.NotEmpty(results)
	require.Len(results, 5)

	for _, result := range results {
		require.EqualError(ErrIncompleteData, result.Err.Error())
	}
}

func (suite *TestBudgetSuite) TestChangesErrorAction() {
	require := require.New(suite.T())
	changes := []Change{
		{
			ID:      uint(1),
			Section: domain.Main,
			Action:  shared.Update,
			Detail: map[string]any{
				"name": "Test 1",
			},
		},
		{
			ID:      uint(1),
			Section: domain.Available,
			Action:  shared.Update,
			Detail: map[string]any{
				"name": "Test 1",
			},
		},
		{
			ID:      uint(2),
			Section: domain.Available,
			Action:  shared.Delete,
			Detail: map[string]any{
				"name": "Test 2",
			},
		},
		{
			ID:      uint(1),
			Section: domain.Bill,
			Action:  shared.Update,
			Detail: map[string]any{
				"description": "Test 1",
			},
		},
		{
			ID:      uint(2),
			Section: domain.Bill,
			Action:  shared.Delete,
			Detail: map[string]any{
				"description": "Test 2",
			},
		},
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	}).Times(5)
	suite.mockBudgetRepo.On("WithTransaction", nil).Return(suite.mockBudgetRepo)
	suite.mockBudgetRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), gorm.ErrRecordNotFound)
	suite.mockBudgetAvailableRepo.On("WithTransaction", nil).Return(suite.mockBudgetAvailableRepo).Times(2)
	suite.mockBudgetAvailableRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), gorm.ErrRecordNotFound)
	suite.mockBudgetAvailableRepo.On("Delete", suite.ctx, mock.Anything).Return(gorm.ErrRecordNotFound)
	suite.mockBudgetBillRepo.On("WithTransaction", nil).Return(suite.mockBudgetBillRepo).Times(2)
	suite.mockBudgetBillRepo.On("Save", suite.ctx, mock.Anything).Return(uint(0), gorm.ErrRecordNotFound)
	suite.mockBudgetBillRepo.On("Delete", suite.ctx, mock.Anything).Return(gorm.ErrRecordNotFound)

	results := suite.app.Changes(suite.ctx, suite.budgetId, changes)

	require.NotEmpty(results)
	require.Len(results, 5)

	for _, result := range results {
		require.EqualError(gorm.ErrRecordNotFound, result.Err.Error())
	}
}

func (suite *TestBudgetSuite) TestChangesErrorInvalidDataTypeMain() {
	require := require.New(suite.T())
	changes := []Change{
		{
			ID:      uint(1),
			Section: domain.Main,
			Action:  shared.Update,
			Detail: map[string]any{
				"year": map[string]any{
					"invalid": 1,
				},
			},
		},
		{
			ID:      uint(2),
			Section: domain.Main,
			Action:  shared.Update,
			Detail: map[string]any{
				"fixedIncome": map[string]any{
					"invalid": 1,
				},
			},
		},
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	}).Times(2)

	results := suite.app.Changes(suite.ctx, suite.budgetId, changes)

	require.NotEmpty(results)
	require.Len(results, 2)

	for _, result := range results {
		require.EqualError(convert.ErrValueIncompatibleType, result.Err.Error())
	}
}

func (suite *TestBudgetSuite) TestChangesErrorInvalidDataTypeAvailable() {
	require := require.New(suite.T())
	changes := []Change{
		{
			ID:      uint(1),
			Section: domain.Available,
			Action:  shared.Update,
			Detail: map[string]any{
				"amount": map[string]any{
					"invalid": 1,
				},
			},
		},
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	})

	results := suite.app.Changes(suite.ctx, suite.budgetId, changes)

	require.NotEmpty(results)
	require.Len(results, 1)

	for _, result := range results {
		require.EqualError(convert.ErrValueIncompatibleType, result.Err.Error())
	}
}

func (suite *TestBudgetSuite) TestChangesErrorInvalidDataTypeBill() {
	require := require.New(suite.T())
	changes := []Change{
		{
			ID:      uint(1),
			Section: domain.Bill,
			Action:  shared.Update,
			Detail: map[string]any{
				"amount": map[string]any{
					"invalid": 1,
				},
			},
		},
		{
			ID:      uint(2),
			Section: domain.Bill,
			Action:  shared.Update,
			Detail: map[string]any{
				"dueDate": map[string]any{
					"invalid": 1,
				},
			},
		},
		{
			ID:      uint(3),
			Section: domain.Bill,
			Action:  shared.Update,
			Detail: map[string]any{
				"complete": map[string]any{
					"invalid": false,
				},
			},
		},
	}
	suite.mockTransactionManager.On("Transaction", mock.AnythingOfType("func(persistent.Transaction) error")).Return(func(fc func(persistent.Transaction) error) error {
		return fc(nil)
	}).Times(3)

	results := suite.app.Changes(suite.ctx, suite.budgetId, changes)

	require.NotEmpty(results)
	require.Len(results, 3)

	for _, result := range results {
		require.EqualError(convert.ErrValueIncompatibleType, result.Err.Error())
	}
}

func (suite *TestBudgetSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockBudgetRepo.On("Delete", suite.ctx, *budgetExpected.ID).Return(nil)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.NoError(err)
}

func (suite *TestBudgetSuite) TestDeleteErrorSearch() {
	require := require.New(suite.T())
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(nil, errExpected)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.EqualError(errExpected, err.Error())
}

func (suite *TestBudgetSuite) TestDeleteError() {
	require := require.New(suite.T())
	name := "Test"
	year := uint16(1)
	month := uint8(1)
	budgetExpected := &domain.Budget{
		ID:     &suite.budgetId,
		Name:   &name,
		Year:   &year,
		Month:  &month,
		UserId: &suite.userId,
	}
	errExpected := errors.New("Error find budget by id")
	suite.mockBudgetRepo.On("Search", suite.ctx, suite.budgetId).Return(budgetExpected, nil)
	suite.mockBudgetRepo.On("Delete", suite.ctx, *budgetExpected.ID).Return(errExpected)

	err := suite.app.Delete(suite.ctx, suite.budgetId)

	require.EqualError(errExpected, err.Error())
}

func TestTestBudgetSuite(t *testing.T) {
	suite.Run(t, new(TestBudgetSuite))
}

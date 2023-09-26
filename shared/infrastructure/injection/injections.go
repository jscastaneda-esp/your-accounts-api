package injection

import (
	budgets_app "your-accounts-api/budgets/application"
	"your-accounts-api/budgets/infrastructure/db/repository/budget"
	"your-accounts-api/budgets/infrastructure/db/repository/budget_available"
	"your-accounts-api/budgets/infrastructure/db/repository/budget_bill"
	logs_app "your-accounts-api/shared/application"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/db/repository/log"
	users_app "your-accounts-api/users/application"
	"your-accounts-api/users/infrastructure/db/repository/user"
	"your-accounts-api/users/infrastructure/db/repository/user_token"
)

var (
	UserApp            users_app.IUserApp
	LogApp             logs_app.ILogApp
	BudgetApp          budgets_app.IBudgetApp
	BudgetAvailableApp budgets_app.IBudgetAvailableApp
)

func LoadInstances() {
	// Repositories
	userRepo := user.NewRepository(db.DB)
	userTokenRepo := user_token.NewRepository(db.DB)
	logRepo := log.NewRepository(db.DB)
	budgetRepo := budget.NewRepository(db.DB)
	budgetAvailableRepo := budget_available.NewRepository(db.DB)
	budgetBillRepo := budget_bill.NewRepository(db.DB)

	// Apps
	UserApp = users_app.NewUserApp(db.Tm, userRepo, userTokenRepo)
	LogApp = logs_app.NewLogApp(db.Tm, logRepo)
	BudgetApp = budgets_app.NewBudgetApp(db.Tm, budgetRepo, budgetAvailableRepo, budgetBillRepo, LogApp)
	BudgetAvailableApp = budgets_app.NewBudgetAvailableApp(db.Tm, budgetAvailableRepo, LogApp)
}

package injection

import (
	budgets_app "your-accounts-api/budgets/application"
	budgets_dom "your-accounts-api/budgets/domain"
	"your-accounts-api/budgets/infrastructure/db/repository/budget"
	"your-accounts-api/budgets/infrastructure/db/repository/budget_available"
	logs_app "your-accounts-api/shared/application"
	logs_dom "your-accounts-api/shared/domain"
	"your-accounts-api/shared/infrastructure/db"
	"your-accounts-api/shared/infrastructure/db/repository/log"
	users_app "your-accounts-api/users/application"
	users_dom "your-accounts-api/users/domain"
	"your-accounts-api/users/infrastructure/db/repository/user"
	"your-accounts-api/users/infrastructure/db/repository/user_token"
)

var (
	UserRepo            users_dom.UserRepository
	UserTokenRepo       users_dom.UserTokenRepository
	LogRepo             logs_dom.LogRepository
	BudgetRepo          budgets_dom.BudgetRepository
	BudgetAvailableRepo budgets_dom.BudgetAvailableRepository
	UserApp             users_app.IUserApp
	LogApp              logs_app.ILogApp
	BudgetApp           budgets_app.IBudgetApp
	BudgetAvailableApp  budgets_app.IBudgetAvailableApp
)

func LoadInstances() {
	// Repositories
	UserRepo = user.NewRepository(db.DB)
	UserTokenRepo = user_token.NewRepository(db.DB)
	LogRepo = log.NewRepository(db.DB)
	BudgetRepo = budget.NewRepository(db.DB)
	BudgetAvailableRepo = budget_available.NewRepository(db.DB)

	// Apps
	UserApp = users_app.NewUserApp(db.Tm, UserRepo, UserTokenRepo)
	LogApp = logs_app.NewLogApp(db.Tm, LogRepo)
	BudgetApp = budgets_app.NewBudgetApp(db.Tm, BudgetRepo, LogApp)
	BudgetAvailableApp = budgets_app.NewBudgetAvailableApp(db.Tm, BudgetAvailableRepo, BudgetApp, LogApp)
}

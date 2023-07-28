package injection

import (
	budget_app "your-accounts-api/budget/application"
	budget_dom "your-accounts-api/budget/domain"
	"your-accounts-api/budget/infrastructure/repository/budget"
	"your-accounts-api/budget/infrastructure/repository/budget_available_balance"
	project_app "your-accounts-api/project/application"
	project_dom "your-accounts-api/project/domain"
	"your-accounts-api/project/infrastructure/repository/project"
	"your-accounts-api/project/infrastructure/repository/project_log"
	"your-accounts-api/shared/infrastructure/db"
	user_app "your-accounts-api/user/application"
	user_dom "your-accounts-api/user/domain"
	"your-accounts-api/user/infrastructure/repository/user"
	"your-accounts-api/user/infrastructure/repository/user_token"
)

var (
	UserRepo                   user_dom.UserRepository
	UserTokenRepo              user_dom.UserTokenRepository
	ProjectRepo                project_dom.ProjectRepository
	ProjectLogRepo             project_dom.ProjectLogRepository
	BudgetRepo                 budget_dom.BudgetRepository
	BudgetAvailableBalanceRepo budget_dom.BudgetAvailableBalanceRepository
	UserApp                    user_app.IUserApp
	ProjectApp                 project_app.IProjectApp
	BudgetApp                  budget_app.IBudgetApp
	BudgetAvailableBalanceApp  budget_app.IBudgetAvailableBalanceApp
)

func LoadInstances() {
	// Repositories
	UserRepo = user.NewRepository(db.DB)
	UserTokenRepo = user_token.NewRepository(db.DB)
	ProjectRepo = project.NewRepository(db.DB)
	ProjectLogRepo = project_log.NewRepository(db.DB)
	BudgetRepo = budget.NewRepository(db.DB)
	BudgetAvailableBalanceRepo = budget_available_balance.NewRepository(db.DB)

	// Apps
	UserApp = user_app.NewUserApp(db.Tm, UserRepo, UserTokenRepo)
	ProjectApp = project_app.NewProjectApp(db.Tm, ProjectRepo, ProjectLogRepo)
	BudgetApp = budget_app.NewBudgetApp(db.Tm, BudgetRepo, ProjectApp)
	BudgetAvailableBalanceApp = budget_app.NewBudgetAvailableBalanceApp(db.Tm, BudgetAvailableBalanceRepo, BudgetApp, ProjectApp)
}

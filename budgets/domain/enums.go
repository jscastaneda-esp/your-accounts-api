package domain

type BudgetBillCategory string

const (
	House                  BudgetBillCategory = "house"
	Entertainment          BudgetBillCategory = "entertainment"
	Personal               BudgetBillCategory = "personal"
	Vehicle_Transportation BudgetBillCategory = "vehicle_transportation"
	Education              BudgetBillCategory = "education"
	Services               BudgetBillCategory = "services"
	Financial              BudgetBillCategory = "financial"
	Saving                 BudgetBillCategory = "saving"
	Others                 BudgetBillCategory = "others"
)

type BudgetSection string

const (
	Main      BudgetSection = "main"
	Available BudgetSection = "available"
	Bill      BudgetSection = "bill"
)

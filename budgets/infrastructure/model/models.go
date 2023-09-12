package model

import (
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/infrastructure/model"
)

type CreateRequest struct {
	Name    string `json:"name" validate:"required_without=CloneId,omitempty,max=40"`
	CloneId *uint  `json:"cloneId" validate:"omitempty,min=1"`
}

type CreateResponse struct {
	model.IDResponse
}

func NewCreateResponse(id uint) CreateResponse {
	return CreateResponse{
		IDResponse: model.NewIDResponse(id),
	}
}

type ReadBudgetResponse struct {
	model.IDResponse
	model.NameResponse
	Year  uint16 `json:"year"`
	Month uint8  `json:"month"`
}

func NewReadBudgetResponse(id uint, name string, year uint16, month uint8) ReadBudgetResponse {
	return ReadBudgetResponse{
		IDResponse:   model.NewIDResponse(id),
		NameResponse: model.NewNameResponse(name),
		Year:         year,
		Month:        month,
	}
}

type ReadResponse struct {
	ReadBudgetResponse
	TotalAvailable float64 `json:"totalAvailable"`
	TotalPending   float64 `json:"totalPending"`
	TotalSaving    float64 `json:"totalSaving"`
	PendingBills   uint8   `json:"pendingBills"`
}

func NewReadResponse(budget domain.Budget) ReadResponse {
	return ReadResponse{
		ReadBudgetResponse: NewReadBudgetResponse(*budget.ID, *budget.Name, *budget.Year, *budget.Month),
		TotalAvailable:     *budget.TotalAvailable,
		TotalPending:       *budget.TotalPending,
		TotalSaving:        *budget.TotalSaving,
		PendingBills:       *budget.PendingBills,
	}
}

type ReadByIDResponseAvailable struct {
	model.IDResponse
	model.NameResponse
	model.AmountResponse
}

type ReadByIDResponseBill struct {
	model.IDResponse
	model.AmountResponse
	Description string                    `json:"description"`
	Payment     float64                   `json:"payment"`
	DueDate     uint8                     `json:"dueDate"`
	Complete    bool                      `json:"complete"`
	Category    domain.BudgetBillCategory `json:"category"`
}

type ReadByIDResponse struct {
	ReadBudgetResponse
	FixedIncome      float64                     `json:"fixedIncome"`
	AdditionalIncome float64                     `json:"additionalIncome"`
	Availables       []ReadByIDResponseAvailable `json:"availables"`
	Bills            []ReadByIDResponseBill      `json:"bills"`
}

func NewReadByIDResponse(budget *domain.Budget) ReadByIDResponse {
	availables := []ReadByIDResponseAvailable{}
	for _, available := range budget.BudgetAvailables {
		availables = append(availables, ReadByIDResponseAvailable{
			IDResponse:     model.NewIDResponse(*available.ID),
			NameResponse:   model.NewNameResponse(*available.Name),
			AmountResponse: model.NewAmountResponse(*available.Amount),
		})
	}

	bills := []ReadByIDResponseBill{}
	for _, bill := range budget.BudgetBills {
		bills = append(bills, ReadByIDResponseBill{
			IDResponse:     model.NewIDResponse(*bill.ID),
			AmountResponse: model.NewAmountResponse(*bill.Amount),
			Description:    *bill.Description,
			Payment:        *bill.Payment,
			DueDate:        *bill.DueDate,
			Complete:       *bill.Complete,
			Category:       *bill.Category,
		})
	}

	return ReadByIDResponse{
		ReadBudgetResponse: NewReadBudgetResponse(*budget.ID, *budget.Name, *budget.Year, *budget.Month),
		FixedIncome:        *budget.FixedIncome,
		AdditionalIncome:   *budget.AdditionalIncome,
		Availables:         availables,
		Bills:              bills,
	}
}

type CreateAvailableRequest struct {
	Name     string `json:"name" validate:"required,max=40"`
	BudgetId uint   `json:"budgetId" validate:"required,min=1"`
}

type CreateAvailableResponse struct {
	model.IDResponse
}

func NewCreateAvailableResponse(id uint) CreateAvailableResponse {
	return CreateAvailableResponse{
		IDResponse: model.NewIDResponse(id),
	}
}

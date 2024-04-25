package model

import (
	"your-accounts-api/budgets/application"
	"your-accounts-api/budgets/domain"
	shared "your-accounts-api/shared/domain"
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
	PendingBills   uint8   `json:"pendingBills"`
}

func NewReadResponse(budget domain.Budget) ReadResponse {
	return ReadResponse{
		ReadBudgetResponse: NewReadBudgetResponse(*budget.ID, *budget.Name, *budget.Year, *budget.Month),
		TotalAvailable:     *budget.TotalAvailable,
		TotalPending:       *budget.TotalPending,
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

func NewReadByIDResponse(budget domain.Budget) ReadByIDResponse {
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

type CreateBillRequest struct {
	Description string                    `json:"description" validate:"required,max=200"`
	Category    domain.BudgetBillCategory `json:"category" validate:"required,oneof='house' 'entertainment' 'personal' 'vehicle_transportation' 'education' 'services' 'financial' 'saving' 'others'"`
	BudgetId    uint                      `json:"budgetId" validate:"required,min=1"`
}

type CreateBillResponse struct {
	model.IDResponse
}

func NewCreateBillResponse(id uint) CreateBillResponse {
	return CreateBillResponse{
		IDResponse: model.NewIDResponse(id),
	}
}

type CreateBillTransactionRequest struct {
	Description string  `json:"description" validate:"required,max=500"`
	Amount      float64 `json:"amount" validate:"required"`
	BillId      uint    `json:"billId" validate:"required,min=1"`
}

type ChangeRequest struct {
	ID      uint                 `json:"id" validate:"required,min=1"`
	Section domain.BudgetSection `json:"section" validate:"required,oneof='main' 'available' 'bill'"`
	Action  shared.Action        `json:"action" validate:"required,oneof='update' 'delete'"`
	Detail  map[string]any       `json:"detail"`
}

type ChangesRequest struct {
	Changes []ChangeRequest `json:"changes" validate:"min=1,dive,required"`
}

type ChangeResponse struct {
	Change ChangeRequest `json:"change"`
	Error  string        `json:"error"`
}

type ChangesResponse struct {
	Changes []ChangeResponse
}

func NewChangeResponse(change application.Change, err string) ChangeResponse {
	return ChangeResponse{
		Change: ChangeRequest{
			ID:      change.ID,
			Section: change.Section,
			Action:  change.Action,
			Detail:  change.Detail,
		},
		Error: err,
	}
}

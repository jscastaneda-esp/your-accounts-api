package model

import (
	"your-accounts-api/budget/domain"
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

type ReadResponse struct {
	model.IDResponse
	Name                  string  `json:"name"`
	Year                  uint16  `json:"year"`
	Month                 uint8   `json:"month"`
	TotalAvailableBalance float64 `json:"totalAvailableBalance"`
	TotalPendingPayment   float64 `json:"totalPendingPayment"`
	TotalBalance          float64 `json:"totalBalance"`
	PendingBills          uint8   `json:"pendingBills"`
}

func NewReadResponse(budget *domain.Budget) ReadResponse {
	return ReadResponse{
		IDResponse:            model.NewIDResponse(*budget.ID),
		Name:                  *budget.Name,
		Year:                  *budget.Year,
		Month:                 *budget.Month,
		TotalAvailableBalance: *budget.TotalAvailableBalance,
		TotalPendingPayment:   *budget.TotalPendingPayment,
		TotalBalance:          *budget.TotalBalance,
		PendingBills:          *budget.PendingBills,
	}
}

type ReadByIDResponse struct {
	model.IDResponse
	Name             string  `json:"name"`
	Year             uint16  `json:"year"`
	Month            uint8   `json:"month"`
	FixedIncome      float64 `json:"fixedIncome"`
	AdditionalIncome float64 `json:"additionalIncome"`
	TotalBalance     float64 `json:"totalBalance"`
	Total            float64 `json:"total"`
	EstimatedBalance float64 `json:"estimatedBalance"`
	ProjectId        uint    `json:"projectId"`
}

func NewReadByIDResponse(budget *domain.Budget) ReadByIDResponse {
	return ReadByIDResponse{
		IDResponse:       model.NewIDResponse(*budget.ID),
		Name:             *budget.Name,
		Year:             *budget.Year,
		Month:            *budget.Month,
		FixedIncome:      *budget.FixedIncome,
		AdditionalIncome: *budget.AdditionalIncome,
		TotalBalance:     *budget.TotalBalance,
		ProjectId:        *budget.ProjectId,
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

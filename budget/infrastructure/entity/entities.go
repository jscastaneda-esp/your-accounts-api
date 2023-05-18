package entity

import (
	"api-your-accounts/budget/domain"
	"api-your-accounts/shared/infrastructure/db/entity"
)

type Budget struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Name                    string                   `gorm:"not null;size:40;uniqueIndex:unq_budget"`
	Year                    uint16                   `gorm:"not null"`
	Month                   uint8                    `gorm:"not null"`
	FixedIncome             float64                  `gorm:"not null;default:0"`
	AdditionalIncome        float64                  `gorm:"not null;default:0"`
	TotalPendingPayment     float64                  `gorm:"not null;default:0"`
	TotalAvailableBalance   float64                  `gorm:"not null;default:0"`
	PendingBills            uint8                    `gorm:"not null;default:0"`
	TotalBalance            float64                  `gorm:"not null;default:0"`
	Total                   float64                  `gorm:"not null;default:0"`
	EstimatedBalance        float64                  `gorm:"not null;default:0"`
	TotalPayment            float64                  `gorm:"not null;default:0"`
	ProjectId               uint                     `gorm:"not null;uniqueIndex:unq_budget"`
	BudgetAvailableBalances []BudgetAvailableBalance `gorm:"foreignKey:BudgetId"`
	BudgetBills             []BudgetBill             `gorm:"foreignKey:BudgetId"`
}

type BudgetAvailableBalance struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Name     string  `gorm:"not null;size:40;uniqueIndex:unq_available"`
	Amount   float64 `gorm:"not null;default:0"`
	BudgetId uint    `gorm:"not null;uniqueIndex:unq_available"`
}

type BudgetBill struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Description            string                    `gorm:"not null;size:200;uniqueIndex:unq_bill"`
	Amount                 float64                   `gorm:"not null;default:0"`
	Payment                float64                   `gorm:"not null;default:0"`
	Shared                 bool                      `gorm:"not null;default:false"`
	DueDate                uint8                     `gorm:"not null;default:0"`
	Complete               bool                      `gorm:"not null;default:false"`
	BudgetId               uint                      `gorm:"not null;uniqueIndex:unq_bill"`
	Category               domain.BudgetBillCategory `gorm:"not null;type:budget_bill_category"`
	BudgetBillTransactions []BudgetBillTransaction   `gorm:"foreignKey:BudgetBillId"`
	BudgetBillShareds      []BudgetBillShared        `gorm:"foreignKey:BudgetBillId"`
}

type BudgetBillTransaction struct {
	entity.BaseModel
	Description  string  `gorm:"not null;size:100"`
	Amount       float64 `gorm:"not null"`
	BudgetBillId uint    `gorm:"not null"`
}

type BudgetBillShared struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Description  string  `gorm:"not null;size:100;uniqueIndex:unq_bill_shared"`
	Amount       float64 `gorm:"not null;default:0"`
	BudgetBillId uint    `gorm:"not null;uniqueIndex:unq_bill_shared"`
}

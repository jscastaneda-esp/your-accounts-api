package infrastructure

import "api-your-accounts/shared/infrastructure/model"

type Budget struct {
	model.BaseModel
	model.BaseUpdateModel
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
	model.BaseModel
	model.BaseUpdateModel
	Name     string  `gorm:"not null;size:40;uniqueIndex:unq_available"`
	Amount   float64 `gorm:"not null;default:0"`
	BudgetId uint    `gorm:"not null;uniqueIndex:unq_available"`
}

type CategoryBill struct {
	model.BaseModel
	model.BaseUpdateModel
	Type        string       `gorm:"not null;size:50;unique"`
	BudgetBills []BudgetBill `gorm:"foreignKey:CategoryId"`
}

type BudgetBill struct {
	model.BaseModel
	model.BaseUpdateModel
	Description            string                  `gorm:"not null;size:200;uniqueIndex:unq_bill"`
	Amount                 float64                 `gorm:"not null;default:0"`
	Payment                float64                 `gorm:"not null;default:0"`
	Shared                 bool                    `gorm:"not null;default:false"`
	DueDate                uint8                   `gorm:"not null;default:0"`
	Complete               bool                    `gorm:"not null;default:false"`
	BudgetId               uint                    `gorm:"not null;uniqueIndex:unq_bill"`
	CategoryId             uint                    `gorm:"not null"`
	BudgetBillTransactions []BudgetBillTransaction `gorm:"foreignKey:BudgetBillId"`
	BudgetBillShareds      []BudgetBillShared      `gorm:"foreignKey:BudgetBillId"`
}

type BudgetBillTransaction struct {
	model.BaseModel
	Description  string  `gorm:"not null;size:100"`
	Amount       float64 `gorm:"not null"`
	BudgetBillId uint    `gorm:"not null"`
}

type BudgetBillShared struct {
	model.BaseModel
	model.BaseUpdateModel
	Description  string  `gorm:"not null;size:100;uniqueIndex:unq_bill_shared"`
	Amount       float64 `gorm:"not null;default:0"`
	BudgetBillId uint    `gorm:"not null;uniqueIndex:unq_bill_shared"`
}

package entity

import (
	"your-accounts-api/budgets/domain"
	"your-accounts-api/shared/infrastructure/db/entity"
)

type Budget struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Name             string            `gorm:"not null;size:40"`
	Year             uint16            `gorm:"not null"`
	Month            uint8             `gorm:"not null"`
	FixedIncome      float64           `gorm:"not null;default:0"`
	AdditionalIncome float64           `gorm:"not null;default:0"`
	TotalPending     float64           `gorm:"not null;default:0"`
	TotalAvailable   float64           `gorm:"not null;default:0"`
	PendingBills     uint8             `gorm:"not null;default:0"`
	UserId           uint              `gorm:"not null"`
	BudgetAvailables []BudgetAvailable `gorm:"foreignKey:BudgetId"`
	BudgetBills      []BudgetBill      `gorm:"foreignKey:BudgetId"`
}

type BudgetAvailable struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Name     string  `gorm:"not null;size:40"`
	Amount   float64 `gorm:"not null;default:0"`
	BudgetId uint    `gorm:"not null"`
}

type BudgetBill struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Description string                    `gorm:"not null;size:200"`
	Amount      float64                   `gorm:"not null;default:0"`
	Payment     float64                   `gorm:"not null;default:0"`
	DueDate     uint8                     `gorm:"not null;default:0"`
	Complete    bool                      `gorm:"not null;default:false"`
	BudgetId    uint                      `gorm:"not null"`
	Category    domain.BudgetBillCategory `gorm:"not null"`
}

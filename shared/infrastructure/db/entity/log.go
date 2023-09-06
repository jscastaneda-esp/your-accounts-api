package entity

import "your-accounts-api/shared/domain"

type Log struct {
	BaseModel
	Description string         `gorm:"not null"`
	Detail      *string        `gorm:"type:json"`
	Code        domain.CodeLog `gorm:"not null;type:enum('budget', 'budget_bill')"`
	ResourceId  uint           `gorm:"not null"`
}

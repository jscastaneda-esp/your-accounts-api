package entity

import "your-accounts-api/shared/domain"

type Log struct {
	BaseModel
	Description string         `gorm:"not null"`
	Detail      map[string]any `gorm:"not null;type:json;serializer:json"`
	Code        domain.CodeLog `gorm:"not null;type:enum('budget', 'budget_bill')"`
	ResourceId  uint           `gorm:"not null"`
}

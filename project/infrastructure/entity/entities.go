package entity

import (
	budget "your-accounts-api/budget/infrastructure/entity"
	"your-accounts-api/project/domain"
	"your-accounts-api/shared/infrastructure/db/entity"
)

type Project struct {
	entity.BaseModel
	entity.BaseUpdateModel
	UserId      uint               `gorm:"not null"`
	Type        domain.ProjectType `gorm:"not null;type:enum('budget')"`
	ProjectLogs []ProjectLog       `gorm:"foreignKey:ProjectId"`
	Budget      budget.Budget      `gorm:"foreignKey:ProjectId"`
}

type ProjectLog struct {
	entity.BaseModel
	Description string  `gorm:"not null"`
	Detail      *string `gorm:"type:json"`
	ProjectId   uint    `gorm:"not null"`
}

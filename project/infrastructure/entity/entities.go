package entity

import (
	budget "api-your-accounts/budget/infrastructure/entity"
	"api-your-accounts/project/domain"
	"api-your-accounts/shared/infrastructure/db/entity"
)

type Project struct {
	entity.BaseModel
	entity.BaseUpdateModel
	UserId      uint               `gorm:"not null"`
	Type        domain.ProjectType `gorm:"not null;type:project_type"`
	ProjectLogs []ProjectLog       `gorm:"foreignKey:ProjectId"`
	Budget      budget.Budget      `gorm:"foreignKey:ProjectId"`
}

type ProjectLog struct {
	entity.BaseModel
	Description string  `gorm:"not null"`
	Detail      *string `gorm:"type:jsonb"`
	ProjectId   uint    `gorm:"not null"`
}

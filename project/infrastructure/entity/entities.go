package entity

import (
	budget "api-your-accounts/budget/infrastructure/entity"
	"api-your-accounts/project/domain"
	"api-your-accounts/shared/infrastructure/db/entity"
)

type Project struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Name       string             `gorm:"not null;size:20;uniqueIndex:unq_project"`
	UserId     uint               `gorm:"not null;uniqueIndex:unq_project"`
	Type       domain.ProjectType `gorm:"not null;size:10;uniqueIndex:unq_project"`
	ProjectLog ProjectLog         `gorm:"foreignKey:ProjectId"`
	Budget     budget.Budget      `gorm:"foreignKey:ProjectId"`
}

type ProjectLog struct {
	entity.BaseModel
	Description string `gorm:"not null"`
	Detail      string `gorm:"type:jsonb;default:'{}';not null"`
	ProjectId   uint   `gorm:"not null"`
}

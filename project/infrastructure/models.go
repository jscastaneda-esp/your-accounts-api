package infrastructure

import (
	"api-your-accounts/budget/infrastructure"
	"api-your-accounts/shared/infrastructure/model"
)

type Project struct {
	model.BaseModel
	model.BaseUpdateModel
	Name   string                `gorm:"not null;size:20;uniqueIndex:unq_project"`
	UserId string                `gorm:"not null;size:32;uniqueIndex:unq_project"`
	Type   string                `gorm:"not null;size:10;uniqueIndex:unq_project"`
	Budget infrastructure.Budget `gorm:"foreignKey:ProjectId"`
}

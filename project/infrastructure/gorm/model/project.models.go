package model

import (
	"api-your-accounts/budget/infrastructure/gorm/model"
	shared "api-your-accounts/shared/infrastructure/gorm/model"
)

type Project struct {
	shared.BaseModel
	shared.BaseUpdateModel
	Name   string       `gorm:"not null;size:20;uniqueIndex:unq_project"`
	UserId string       `gorm:"not null;size:32;uniqueIndex:unq_project"`
	Type   string       `gorm:"not null;size:10;uniqueIndex:unq_project"`
	Budget model.Budget `gorm:"foreignKey:ProjectId"`
}

package entity

import (
	budget "api-your-accounts/budget/infrastructure/entity"
	"api-your-accounts/shared/infrastructure/db/entity"
)

type Project struct {
	entity.BaseModel
	entity.BaseUpdateModel
	Name   string        `gorm:"not null;size:20;uniqueIndex:unq_project"`
	UserId string        `gorm:"not null;size:32;uniqueIndex:unq_project"`
	Type   string        `gorm:"not null;size:10;uniqueIndex:unq_project"`
	Budget budget.Budget `gorm:"foreignKey:ProjectId"`
}

type ProjectLog struct {
	entity.MongoBaseModel
	Description string
	Detail      map[string]any `bson:"inline"`
	ProjectId   uint           `bson:"project_id"`
}

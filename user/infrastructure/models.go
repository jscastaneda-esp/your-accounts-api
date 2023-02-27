package infrastructure

import (
	"api-your-accounts/project/infrastructure"
	"api-your-accounts/shared/infrastructure/model"
)

type User struct {
	model.BaseModel
	model.BaseUpdateModel
	UUID     string                   `gorm:"not null;size:32;unique"`
	Email    string                   `gorm:"not null;unique"`
	Projects []infrastructure.Project `gorm:"foreignKey:UserId;references:UUID"`
}

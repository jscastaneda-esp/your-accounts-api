package entity

import (
	project "api-your-accounts/project/infrastructure/entity"
	"api-your-accounts/shared/infrastructure/db/entity"
	"time"
)

type User struct {
	entity.BaseModel
	entity.BaseUpdateModel
	UUID       string            `gorm:"not null;size:32;unique"`
	Email      string            `gorm:"not null;unique"`
	Projects   []project.Project `gorm:"foreignKey:UserId"`
	UserTokens []UserToken       `gorm:"foreignKey:UserId"`
}

type UserToken struct {
	entity.BaseModel
	Token          string `gorm:"not null;size:2000;uniqueIndex:unq_token"`
	UserId         uint   `gorm:"not null;uniqueIndex:unq_token"`
	RefreshedBy    *uint
	ExpiresAt      time.Time `gorm:"not null"`
	RefreshedAt    *time.Time
	RefreshedToken *UserToken `gorm:"foreignKey:RefreshedBy"`
}

package entity

import (
	"time"
	project "your-accounts-api/project/infrastructure/entity"
	"your-accounts-api/shared/infrastructure/db/entity"
)

type User struct {
	entity.BaseModel
	entity.BaseUpdateModel
	UID        string            `gorm:"not null;size:32;unique"`
	Email      string            `gorm:"not null;unique"`
	Projects   []project.Project `gorm:"foreignKey:UserId"`
	UserTokens []UserToken       `gorm:"foreignKey:UserId"`
}

type UserToken struct {
	entity.BaseModel
	Token     string    `gorm:"not null;size:2000"`
	UserId    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

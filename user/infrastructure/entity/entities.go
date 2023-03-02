package entity

import (
	project "api-your-accounts/project/infrastructure"
	"api-your-accounts/shared/infrastructure/db/entity"
	"time"
)

type User struct {
	entity.BaseModel
	entity.BaseUpdateModel
	UUID     string            `gorm:"not null;size:32;unique"`
	Email    string            `gorm:"not null;unique"`
	Projects []project.Project `gorm:"foreignKey:UserId;references:UUID"`
}

type SessionLog struct {
	entity.MongoBaseModel
	Description string
	Detail      map[string]interface{} `bson:"inline"`
	UserId      string                 `bson:"user_id"`
	EndedAt     time.Time              `bson:"ended_at"`
}

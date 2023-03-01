package infrastructure

import (
	"api-your-accounts/project/infrastructure"
	"api-your-accounts/shared/infrastructure/model"
	"time"
)

type User struct {
	model.BaseModel
	model.BaseUpdateModel
	UUID     string                   `gorm:"not null;size:32;unique"`
	Email    string                   `gorm:"not null;unique"`
	Projects []infrastructure.Project `gorm:"foreignKey:UserId;references:UUID"`
}

type SessionLog struct {
	model.MongoBaseModel
	Description string
	Detail      map[string]interface{} `bson:"inline"`
	UserId      string                 `bson:"user_id"`
	EndedAt     time.Time              `bson:"ended_at"`
}

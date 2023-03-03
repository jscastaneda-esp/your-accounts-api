package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

type BaseUpdateModel struct {
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
}

type MongoBaseModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

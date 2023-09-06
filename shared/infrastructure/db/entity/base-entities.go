package entity

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime:true"`
}

type BaseUpdateModel struct {
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
}

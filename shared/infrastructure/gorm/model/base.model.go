package model

import "time"

type BaseModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

type BaseUpdateModel struct {
	UpdatedAt time.Time `gorm:"autoUpdateTime:true"`
}

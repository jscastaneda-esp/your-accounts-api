package entity

type Log struct {
	BaseModel
	Description string  `gorm:"not null"`
	Detail      *string `gorm:"type:json"`
	ResourceId  uint    `gorm:"not null"`
}

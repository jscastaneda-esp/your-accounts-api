package domain

import (
	"time"
)

type UserToken struct {
	ID          uint
	Token       string
	UserId      uint
	RefreshedId uint
	CreatedAt   time.Time
	ExpiresAt   time.Time
	RefreshedAt time.Time
}

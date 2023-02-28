package domain

import (
	"context"
	"time"
)

type User struct {
	Id        uint
	UUID      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
}

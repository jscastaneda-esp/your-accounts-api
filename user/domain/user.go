package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint
	UUID      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*User, error)
	ExistsByUUID(ctx context.Context, uuid string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
}

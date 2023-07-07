package domain

import (
	"context"
	"time"
	"your-accounts-api/shared/domain/persistent"
)

type User struct {
	ID        uint
	UUID      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//go:generate mockery --name UserRepository --filename user-repository.go
type UserRepository interface {
	persistent.TransactionRepository[UserRepository]
	persistent.CreateRepository[User]
	FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*User, error)
	ExistsByUUID(ctx context.Context, uuid string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

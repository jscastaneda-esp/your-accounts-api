package domain

import (
	"api-your-accounts/shared/domain/transaction"
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

//go:generate mockery --name UserRepository --filename user-repository.go
type UserRepository interface {
	WithTransaction(tx transaction.Transaction) UserRepository
	FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*User, error)
	ExistsByUUID(ctx context.Context, uuid string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *User) (*User, error)
}

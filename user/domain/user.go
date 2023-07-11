package domain

import (
	"context"
	"time"
	"your-accounts-api/shared/domain/persistent"
)

type User struct {
	ID        uint
	UID       string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//go:generate mockery --name UserRepository --filename user-repository.go
type UserRepository interface {
	persistent.TransactionRepository[UserRepository]
	persistent.CreateRepository[User]
	FindByUIDAndEmail(ctx context.Context, uid string, email string) (*User, error)
	ExistsByUID(ctx context.Context, uid string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

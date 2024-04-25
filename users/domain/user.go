package domain

import (
	"context"
	"your-accounts-api/shared/domain/persistent"
)

type User struct {
	ID    uint
	Email string
}

type UserRepository interface {
	persistent.TransactionRepository[UserRepository]
	persistent.SaveRepository[User]
	persistent.SearchByExampleRepository[User]
	ExistsByExample(ctx context.Context, example User) (bool, error)
}

package domain

import (
	"context"
	"time"
	"your-accounts-api/shared/domain/persistent"
)

type UserToken struct {
	ID        uint
	Token     string
	UserId    uint
	ExpiresAt time.Time
}

type UserTokenRepository interface {
	persistent.TransactionRepository[UserTokenRepository]
	persistent.SaveRepository[UserToken]
	persistent.SearchByExampleRepository[UserToken]
	DeleteByExpiresAtGreaterThanNow(ctx context.Context) error
}

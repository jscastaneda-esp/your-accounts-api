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
	CreatedAt time.Time
	ExpiresAt time.Time
}

//go:generate mockery --name UserTokenRepository --filename user-token-repository.go
type UserTokenRepository interface {
	persistent.TransactionRepository[UserTokenRepository]
	persistent.CreateRepository[UserToken]
	FindByTokenAndUserId(ctx context.Context, token string, userId uint) (*UserToken, error)
	persistent.UpdateRepository[UserToken]
}

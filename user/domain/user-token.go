package domain

import (
	"context"
	"time"
)

type UserToken struct {
	ID          uint
	Token       string
	UserId      uint
	RefreshedId *uint
	CreatedAt   time.Time
	ExpiresAt   time.Time
	RefreshedAt *time.Time
}

//go:generate mockery --name UserTokenRepository --filename user-token-repository.go
type UserTokenRepository interface {
	Create(ctx context.Context, userToken *UserToken) (*UserToken, error)
}
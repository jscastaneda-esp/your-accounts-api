package domain

import (
	"api-your-accounts/shared/domain/persistent"
	"context"
	"time"
)

type Project struct {
	ID        uint
	UserId    uint
	Type      ProjectType
	CreatedAt time.Time
	UpdatedAt time.Time
}

//go:generate mockery --name ProjectRepository --filename project-repository.go
type ProjectRepository interface {
	persistent.TransactionRepository[ProjectRepository]
	persistent.CreateRepository[Project]
	persistent.ReadRepository[Project, uint]
	FindByUserId(ctx context.Context, userId uint) ([]*Project, error)
	persistent.DeleteRepository[uint]
}
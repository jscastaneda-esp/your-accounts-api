package domain

import (
	"context"
	"time"
	"your-accounts-api/shared/domain/persistent"
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
	persistent.ReadRepository[Project]
	FindByUserId(ctx context.Context, userId uint) ([]*Project, error)
	persistent.DeleteRepository
}

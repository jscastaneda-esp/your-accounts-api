package domain

import (
	"api-your-accounts/shared/domain/persistent"
	"context"
	"time"
)

type Project struct {
	ID        uint
	Name      string
	UserId    uint
	Type      ProjectType
	CreatedAt time.Time
	UpdatedAt time.Time
}

//go:generate mockery --name ProjectRepository --filename project-repository.go
type ProjectRepository interface {
	persistent.TransactionRepository[Project]
	persistent.CreateRepository[Project]
	FindByUserId(ctx context.Context, userId uint) ([]*Project, error)
	persistent.DeleteRepository[uint]
}

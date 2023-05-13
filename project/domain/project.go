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
	persistent.TransactionRepository[ProjectRepository]
	persistent.CreateRepository[Project]
	FindByUserId(ctx context.Context, userId uint) ([]*Project, error)
	ExistsByNameAndUserIdAndType(ctx context.Context, name string, userId uint, typeP ProjectType) (bool, error)
	persistent.DeleteRepository[uint]
}

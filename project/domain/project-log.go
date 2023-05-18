package domain

import (
	"api-your-accounts/shared/domain/persistent"
	"context"
	"time"
)

type ProjectLog struct {
	ID          uint
	Description string
	Detail      *string
	ProjectId   uint
	CreatedAt   time.Time
}

//go:generate mockery --name ProjectLogRepository --filename project-log-repository.go
type ProjectLogRepository interface {
	persistent.TransactionRepository[ProjectLogRepository]
	persistent.CreateRepository[ProjectLog]
	FindByProjectId(ctx context.Context, projectId uint) ([]*ProjectLog, error)
}

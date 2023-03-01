package domain

import (
	"context"
	"time"
)

type ProjectLog struct {
	ID          string
	Description string
	Detail      map[string]interface{}
	ProjectId   uint
	CreatedAt   time.Time
}

type ProjectLogRepository interface {
	FindByProjectId(ctx context.Context, projectId uint) ([]*ProjectLog, error)
	CreateAll(ctx context.Context, logs []*ProjectLog) error
}

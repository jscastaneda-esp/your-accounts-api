package domain

import (
	"context"
	"time"
	"your-accounts-api/shared/domain/persistent"
)

type Log struct {
	ID          uint
	Description string
	Detail      map[string]any
	Code        CodeLog
	ResourceId  uint
	CreatedAt   time.Time
}

type LogRepository interface {
	persistent.TransactionRepository[LogRepository]
	persistent.SaveRepository[Log]
	persistent.SearchAllByExampleRepository[Log]
	DeleteByResourceIdNotExists(ctx context.Context) error
	SearchResourceIdsWithLimitExceeded(ctx context.Context) ([]uint, error)
	DeleteByResourceIdAndIdLessThanLimit(ctx context.Context, resourceId uint) error
}

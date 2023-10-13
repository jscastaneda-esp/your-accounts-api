package domain

import (
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

//go:generate mockery --name LogRepository --filename log-repository.go
type LogRepository interface {
	persistent.TransactionRepository[LogRepository]
	persistent.SaveRepository[Log]
	persistent.SearchAllByExampleRepository[Log]
}

package domain

import (
	"time"
	"your-accounts-api/shared/domain/persistent"
)

type Log struct {
	ID          uint
	Description string
	Detail      *string
	Code        CodeLog
	ResourceId  uint
	CreatedAt   time.Time
}

//go:generate mockery --name LogRepository --filename log-repository.go
type LogRepository interface {
	persistent.TransactionRepository[LogRepository]
	persistent.SaveRepository[Log]
	persistent.SearchAllByExample[Log]
}

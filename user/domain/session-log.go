package domain

import (
	"context"
	"time"
)

type SessionLog struct {
	ID          string
	Description string
	Detail      map[string]any
	UserId      string
	CreatedAt   time.Time
	EndedAt     time.Time
}

//go:generate mockery --name SessionLogRepository
type SessionLogRepository interface {
	Create(ctx context.Context, session *SessionLog) (*SessionLog, error)
	Update(ctx context.Context, session *SessionLog) (*SessionLog, error)
}

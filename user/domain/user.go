package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint
	UUID      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	FindByUUIDAndEmail(ctx context.Context, uuid string, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
}

type SessionLog struct {
	ID          string
	Description string
	Detail      map[string]interface{}
	UserId      string
	CreatedAt   time.Time
	EndedAt     time.Time
}

type SessionLogRepository interface {
	Create(ctx context.Context, session *SessionLog) (*SessionLog, error)
	Update(ctx context.Context, session *SessionLog) (*SessionLog, error)
}

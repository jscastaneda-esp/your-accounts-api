package model

import (
	"api-your-accounts/project/domain"
	"time"
)

type CreateRequest struct {
	Name    string             `json:"name,omitempty" validate:"required_without=CloneId,omitempty,max=40"`
	Type    domain.ProjectType `json:"type,omitempty" validate:"required_without=CloneId,omitempty,oneof='budget'"`
	UserId  uint               `json:"userId,omitempty" validate:"required_without=CloneId,omitempty,min=1"`
	CloneId *uint              `json:"cloneId,omitempty" validate:"omitempty,min=1"`
}

type CreateResponse struct {
	ID uint `json:"id,omitempty"`
}

type ReadResponse struct {
	ID   uint               `json:"id,omitempty"`
	Name string             `json:"name,omitempty"`
	Type domain.ProjectType `json:"type,omitempty"`
}

type ReadLogsResponse struct {
	ID          uint      `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

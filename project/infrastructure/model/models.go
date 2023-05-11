package model

import (
	"api-your-accounts/project/domain"
	"time"
)

type CreateRequest struct {
	UserID  uint               `json:"userId,omitempty" validate:"required,min=1"`
	Type    domain.ProjectType `json:"type,omitempty" validate:"required,oneof='budget'"`
	CloneID *uint              `json:"cloneId,omitempty" validate:"omitempty,min=1"`
}

type CreateResponse struct {
	ID uint `json:"id,omitempty"`
}

type ReadResponse struct {
	ID   uint               `json:"id,omitempty"`
	Name string             `json:"name,omitempty"`
	Type domain.ProjectType `json:"type,omitempty"`
}

type ReadTransactionResponse struct {
	ID          uint      `json:"id,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

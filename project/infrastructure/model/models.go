package model

import (
	"time"
	"your-accounts-api/project/domain"
	"your-accounts-api/shared/infrastructure/model"
)

type ReadLogsResponse struct {
	model.IDResponse
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

func NewReadLogsResponse(log *domain.ProjectLog) *ReadLogsResponse {
	return &ReadLogsResponse{
		IDResponse: model.IDResponse{
			ID: log.ID,
		},
		Description: log.Description,
		CreatedAt:   log.CreatedAt,
	}
}

package model

import (
	"time"
	"your-accounts-api/shared/domain"
)

type ReadLogsResponse struct {
	IDResponse
	Description string         `json:"description"`
	Detail      map[string]any `json:"detail"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func NewReadLogsResponse(log domain.Log) ReadLogsResponse {
	return ReadLogsResponse{
		IDResponse:  NewIDResponse(log.ID),
		Description: log.Description,
		Detail:      log.Detail,
		CreatedAt:   log.CreatedAt,
	}
}

// TODO: Pendientes tests

package application

import (
	"api-your-accounts/project/domain"
	"context"
)

func FindLogsByProjectId(repo domain.ProjectLogRepository, ctx context.Context, projectId uint) ([]*domain.ProjectLog, error) {
	return nil, nil
}

func CreateLog(repo domain.ProjectLogRepository, ctx context.Context, log *domain.ProjectLog) (*domain.ProjectLog, error) {
	return nil, nil
}

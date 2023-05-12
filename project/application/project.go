package application

import (
	"api-your-accounts/project/domain"
	"api-your-accounts/shared/domain/persistent"
	"context"
)

//go:generate mockery --name IProjectApp --filename project-app.go
type IProjectApp interface {
	Create(ctx context.Context, project *domain.Project) (*domain.Project, error)
	Clone(ctx context.Context, project *domain.Project, baseId uint) (*domain.Project, error)
	FindByUser(ctx context.Context, userId uint) ([]*domain.Project, error)
	FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error)
	Delete(ctx context.Context, id uint) error
}

type projectApp struct {
	tm             persistent.TransactionManager
	projectRepo    domain.ProjectRepository
	projectLogRepo domain.ProjectLogRepository
}

func (app *projectApp) Create(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) Clone(ctx context.Context, project *domain.Project, baseId uint) (*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) FindById(ctx context.Context, id uint) (*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) FindByUser(ctx context.Context, userId uint) ([]*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) Delete(ctx context.Context, id uint) error {
	panic("not implemented") // TODO: Implement
}

func NewProjectApp(tm persistent.TransactionManager, projectRepo domain.ProjectRepository, projectLogRepo domain.ProjectLogRepository) IProjectApp {
	return &projectApp{tm, projectRepo, projectLogRepo}
}

package application

import (
	"api-your-accounts/project/domain"
	"context"
)

//go:generate mockery --name IProjectApp --filename project-app.go
type IProjectApp interface {
	Create(ctx context.Context, project *domain.Project) (*domain.Project, error)
	ReadByUser(ctx context.Context, userId uint) ([]*domain.Project, error)
	Delete(ctx context.Context, id uint) error
}

type projectApp struct {
}

func (app *projectApp) Create(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) ReadByUser(ctx context.Context, userId uint) ([]*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) Delete(ctx context.Context, id uint) error {
	panic("not implemented") // TODO: Implement
}

func NewProjectApp() IProjectApp {
	return &projectApp{}
}

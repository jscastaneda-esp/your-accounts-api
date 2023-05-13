package application

import (
	"api-your-accounts/project/domain"
	"api-your-accounts/shared/domain/persistent"
	"context"
	"errors"
)

var (
	ErrProjectAlreadyExists = errors.New("project already exists")
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
	exists, err := app.projectRepo.ExistsByNameAndUserIdAndType(ctx, project.Name, project.UserId, project.Type)
	if err != nil {
		return nil, err
	} else if exists {
		return nil, ErrProjectAlreadyExists
	}

	var newProject *domain.Project
	err = app.tm.Transaction(func(tx persistent.Transaction) error {
		projectRepo := app.projectRepo.WithTransaction(tx)
		newProject = &domain.Project{
			Name:   project.Name,
			UserId: project.UserId,
			Type:   project.Type,
		}
		var err error
		newProject, err = projectRepo.Create(ctx, newProject)
		if err != nil {
			return err
		}

		projectLogRepo := app.projectLogRepo.WithTransaction(tx)
		newLog := &domain.ProjectLog{
			Description: "Creación",
			ProjectId:   newProject.ID,
		}
		_, err = projectLogRepo.Create(ctx, newLog)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newProject, nil
}

func (app *projectApp) Clone(ctx context.Context, project *domain.Project, baseId uint) (*domain.Project, error) {
	panic("not implemented") // TODO: Implement
}

func (app *projectApp) FindByUser(ctx context.Context, userId uint) ([]*domain.Project, error) {
	projects, err := app.projectRepo.FindByUserId(ctx, userId)
	if err != nil {
		return []*domain.Project{}, nil
	}

	return projects, nil
}

func (app *projectApp) FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error) {
	logs, err := app.projectLogRepo.FindByProjectId(ctx, projectId)
	if err != nil {
		return []*domain.ProjectLog{}, nil
	}

	return logs, nil
}

func (app *projectApp) Delete(ctx context.Context, id uint) error {
	return app.projectRepo.Delete(ctx, id)
}

func NewProjectApp(tm persistent.TransactionManager, projectRepo domain.ProjectRepository, projectLogRepo domain.ProjectLogRepository) IProjectApp {
	return &projectApp{tm, projectRepo, projectLogRepo}
}

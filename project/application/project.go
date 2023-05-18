package application

import (
	"api-your-accounts/project/domain"
	"api-your-accounts/shared/domain/persistent"
	"context"
	"fmt"
)

//go:generate mockery --name IProjectApp --filename project-app.go
type IProjectApp interface {
	Create(ctx context.Context, project *domain.Project, cloneId *uint) (*domain.Project, error)
	Clone(ctx context.Context, baseId uint) (*domain.Project, error)
	FindByUser(ctx context.Context, userId uint) ([]*domain.Project, error)
	FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error)
	Delete(ctx context.Context, id uint) error
}

type projectApp struct {
	tm             persistent.TransactionManager
	projectRepo    domain.ProjectRepository
	projectLogRepo domain.ProjectLogRepository
}

func (app *projectApp) Create(ctx context.Context, project *domain.Project, cloneId *uint) (*domain.Project, error) {
	var newProject *domain.Project
	err := app.tm.Transaction(func(tx persistent.Transaction) error {
		projectRepo := app.projectRepo.WithTransaction(tx)
		newProject = &domain.Project{
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
		if cloneId != nil {
			detail := fmt.Sprintf(`{"cloneId": %d}`, *cloneId)
			newLog.Description += fmt.Sprintf(" a partir de %d", *cloneId)
			newLog.Detail = &detail
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

func (app *projectApp) Clone(ctx context.Context, baseId uint) (*domain.Project, error) {
	baseProject, err := app.projectRepo.FindById(ctx, baseId)
	if err != nil {
		return nil, err
	}

	newProject := &domain.Project{
		UserId: baseProject.UserId,
		Type:   baseProject.Type,
	}
	return app.Create(ctx, newProject, &baseId)
}

func (app *projectApp) FindByUser(ctx context.Context, userId uint) ([]*domain.Project, error) {
	projects, err := app.projectRepo.FindByUserId(ctx, userId)
	if err != nil {
		return []*domain.Project{}, err
	}

	return projects, nil
}

func (app *projectApp) FindLogsByProject(ctx context.Context, projectId uint) ([]*domain.ProjectLog, error) {
	logs, err := app.projectLogRepo.FindByProjectId(ctx, projectId)
	if err != nil {
		return []*domain.ProjectLog{}, err
	}

	return logs, nil
}

func (app *projectApp) Delete(ctx context.Context, id uint) error {
	return app.projectRepo.Delete(ctx, id)
}

func NewProjectApp(tm persistent.TransactionManager, projectRepo domain.ProjectRepository, projectLogRepo domain.ProjectLogRepository) IProjectApp {
	return &projectApp{tm, projectRepo, projectLogRepo}
}

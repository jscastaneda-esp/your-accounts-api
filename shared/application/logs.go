package application

import (
	"context"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
)

//go:generate mockery --name ILogApp --filename log-app.go
type ILogApp interface {
	CreateLog(ctx context.Context, description string, code domain.CodeLog, resourceId uint, detail map[string]any, tx persistent.Transaction) error
	FindLogsByProject(ctx context.Context, code domain.CodeLog, resourceId uint) ([]domain.Log, error)
}

type logApp struct {
	tm      persistent.TransactionManager
	logRepo domain.LogRepository
}

func (app *logApp) CreateLog(ctx context.Context, description string, code domain.CodeLog, resourceId uint, detail map[string]any, tx persistent.Transaction) error {
	projectLogRepo := app.logRepo.WithTransaction(tx)
	newLog := domain.Log{
		Description: description,
		Detail:      detail,
		Code:        code,
		ResourceId:  resourceId,
	}
	_, err := projectLogRepo.Save(ctx, newLog)
	if err != nil {
		return err
	}

	return nil
}

func (app *logApp) FindLogsByProject(ctx context.Context, code domain.CodeLog, resourceId uint) ([]domain.Log, error) {
	example := domain.Log{
		Code:       code,
		ResourceId: resourceId,
	}
	logs, err := app.logRepo.SearchAllByExample(ctx, example)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func NewLogApp(tm persistent.TransactionManager, logRepo domain.LogRepository) ILogApp {
	return &logApp{tm, logRepo}
}

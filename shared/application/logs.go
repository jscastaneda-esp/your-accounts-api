package application

import (
	"context"
	"sync"
	"your-accounts-api/shared/domain"
	"your-accounts-api/shared/domain/persistent"
)

type ILogApp interface {
	Create(ctx context.Context, description string, code domain.CodeLog, resourceId uint, detail map[string]any, tx persistent.Transaction) error
	FindByProject(ctx context.Context, code domain.CodeLog, resourceId uint) ([]domain.Log, error)
	DeleteOrphan(ctx context.Context) error
	DeleteOld(ctx context.Context) error
}

type logApp struct {
	tm      persistent.TransactionManager
	logRepo domain.LogRepository
	mu      sync.Mutex
}

func (app *logApp) Create(ctx context.Context, description string, code domain.CodeLog, resourceId uint, detail map[string]any, tx persistent.Transaction) error {
	projectLogRepo := app.logRepo.WithTransaction(tx)
	newLog := domain.Log{
		Description: description,
		Detail:      detail,
		Code:        code,
		ResourceId:  resourceId,
	}

	app.mu.Lock()
	defer app.mu.Unlock()
	_, err := projectLogRepo.Save(ctx, newLog)
	if err != nil {
		return err
	}

	return nil
}

func (app *logApp) FindByProject(ctx context.Context, code domain.CodeLog, resourceId uint) ([]domain.Log, error) {
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

func (app *logApp) DeleteOrphan(ctx context.Context) error {
	err := app.logRepo.DeleteByResourceIdNotExists(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (app *logApp) DeleteOld(ctx context.Context) error {
	resourceIds, err := app.logRepo.SearchResourceIdsWithLimitExceeded(ctx)
	if err != nil {
		return err
	}

	if len(resourceIds) == 0 {
		return nil
	}

	return app.tm.Transaction(func(tx persistent.Transaction) error {
		logRepo := app.logRepo.WithTransaction(tx)

		for _, resourceId := range resourceIds {
			err := logRepo.DeleteByResourceIdAndIdLessThanLimit(ctx, resourceId)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func NewLogApp(tm persistent.TransactionManager, logRepo domain.LogRepository) ILogApp {
	return &logApp{
		tm:      tm,
		logRepo: logRepo,
	}
}

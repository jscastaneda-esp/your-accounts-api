package domain

import (
	"your-accounts-api/shared/domain/persistent"
)

type Project struct {
	ID     uint
	UserId uint
	Type   ProjectType
}

//go:generate mockery --name ProjectRepository --filename project-repository.go
type ProjectRepository interface {
	persistent.TransactionRepository[ProjectRepository]
	persistent.SaveRepository[Project]
	persistent.DeleteRepository
}

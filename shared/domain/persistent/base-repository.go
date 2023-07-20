package persistent

import "context"

type CreateRepository[D any] interface {
	Create(ctx context.Context, domain D) (uint, error)
}

type ReadRepository[D any] interface {
	FindById(ctx context.Context, id uint) (*D, error)
}

type UpdateRepository[D any] interface {
	Update(ctx context.Context, domain D) error
}

type DeleteRepository interface {
	Delete(ctx context.Context, id uint) error
}

type CrudRepository[D any, ID any] interface {
	CreateRepository[D]
	ReadRepository[D]
	UpdateRepository[D]
	DeleteRepository
}

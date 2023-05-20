package persistent

import "context"

type CreateRepository[D any] interface {
	Create(ctx context.Context, t *D) (*D, error)
}

type ReadRepository[D any, ID any] interface {
	FindById(ctx context.Context, id ID) (*D, error)
}

type UpdateRepository[D any] interface {
	Update(ctx context.Context, userToken *D) (*D, error)
}

type DeleteRepository[ID any] interface {
	Delete(ctx context.Context, id ID) error
}

type CrudRepository[D any, ID any] interface {
	CreateRepository[D]
	ReadRepository[D, ID]
	UpdateRepository[D]
	DeleteRepository[ID]
}

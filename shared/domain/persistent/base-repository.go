package persistent

import "context"

type CreateRepository[T any] interface {
	Create(ctx context.Context, t *T) (*T, error)
}

type UpdateRepository[T any] interface {
	Update(ctx context.Context, userToken *T) (*T, error)
}

type DeleteRepository[T any] interface {
	Delete(ctx context.Context, id T) error
}

type CrudRepository[E any, ID any] interface {
	CreateRepository[E]
	UpdateRepository[E]
	DeleteRepository[ID]
}

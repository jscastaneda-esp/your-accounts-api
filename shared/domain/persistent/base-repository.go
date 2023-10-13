package persistent

import "context"

type SaveRepository[D any] interface {
	Save(ctx context.Context, domain D) (uint, error)
}

type SaveAllRepository[D any] interface {
	SaveAll(ctx context.Context, domains []D) error
}

type SearchRepository[D any] interface {
	Search(ctx context.Context, id uint) (*D, error)
}

type SearchByExampleRepository[D any] interface {
	SearchByExample(ctx context.Context, example D) (*D, error)
}

type SearchAllByExampleRepository[D any] interface {
	SearchAllByExample(ctx context.Context, example D) ([]D, error)
}

type DeleteRepository interface {
	Delete(ctx context.Context, id uint) error
}

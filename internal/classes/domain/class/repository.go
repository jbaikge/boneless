package class

import "context"

type Repository interface {
	AddClass(ctx context.Context, c *Class) error
	GetClass(ctx context.Context, id string) (*Class, error)
	UpdateClass(
		ctx context.Context,
		id string,
		updateFn func(ctx context.Context, c *Class) (*Class, error),
	) error
}

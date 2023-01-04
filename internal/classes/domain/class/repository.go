package class

import "context"

type Repository interface {
	AddClass(ctx context.Context, class *Class) error
	GetClass(ctx context.Context, id string) (*Class, error)
	UpdateClass(ctx context.Context, class *Class) error
}

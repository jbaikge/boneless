package decorator

import "context"

type QueryHandler[Query any, Result any] interface {
	Handle(context.Context, Query) (Result, error)
}

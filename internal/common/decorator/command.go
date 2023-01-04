package decorator

import "context"

type CommandHandler[Command any] interface {
	Handle(ctx context.Context, cmd Command) error
}

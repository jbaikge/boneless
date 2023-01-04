package service

import (
	"context"
	"errors"

	"github.com/jbaikge/boneless/internal/classes/app"
	"github.com/jbaikge/boneless/internal/classes/app/command"
	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/storage"
	"golang.org/x/exp/slog"
)

func NewApplication(ctx context.Context, store storage.Option) app.Application {
	var classRepo class.Repository
	switch store {
	case storage.DynamoDB:
		// TODO
	case storage.SQLite:
		// TODO
	case storage.Memory:
		// TODO
	default:
		slog.Error("unable to init class application", errors.New("invalid storage option"), "option", store)
	}

	return app.Application{
		Commands: app.Commands{
			AddClass: command.NewAddClassHandler(classRepo),
		},
	}
}

package service

import (
	"context"
	"errors"

	"github.com/jbaikge/boneless/internal/classes/adapters"
	"github.com/jbaikge/boneless/internal/classes/app"
	"github.com/jbaikge/boneless/internal/classes/app/command"
	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/logger"
	"github.com/jbaikge/boneless/internal/common/storage"
)

func NewApplication(ctx context.Context, store storage.Option, log logger.Logger) app.Application {
	var classRepo class.Repository
	switch store {
	case storage.DynamoDB:
		// TODO
	case storage.SQLite:
		// TODO
	case storage.Memory:
		classRepo = adapters.NewClassMemoryRepository()
	default:
		log.Error("unable to init class application", errors.New("invalid storage option"), "option", store)
	}

	return app.Application{
		Commands: app.Commands{
			AddClass: command.NewAddClassHandler(classRepo, log),
		},
	}
}

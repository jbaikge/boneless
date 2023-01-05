package command

import (
	"context"

	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/decorator"
	"github.com/jbaikge/boneless/internal/common/logger"
)

type AddClass struct {
	Class *class.Class
}

type AddClassHandler decorator.CommandHandler[AddClass]

type addClassHandler struct {
	classRepo class.Repository
	log       logger.Logger
}

func NewAddClassHandler(classRepo class.Repository, log logger.Logger) AddClassHandler {
	return addClassHandler{
		classRepo: classRepo,
		log:       log,
	}
}

func (h addClassHandler) Handle(ctx context.Context, cmd AddClass) (err error) {
	h.log.Debug("adding class", "cmd", cmd)
	return h.classRepo.AddClass(ctx, cmd.Class)
}

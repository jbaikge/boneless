package command

import (
	"context"

	"github.com/jbaikge/boneless/internal/classes/domain/class"
	"github.com/jbaikge/boneless/internal/common/decorator"
	"golang.org/x/exp/slog"
)

type AddClass struct {
	Class *class.Class
}

type AddClassHandler decorator.CommandHandler[AddClass]

type addClassHandler struct {
	classRepo class.Repository
}

func NewAddClassHandler(classRepo class.Repository) AddClassHandler {
	return addClassHandler{
		classRepo: classRepo,
	}
}

func (h addClassHandler) Handle(ctx context.Context, cmd AddClass) (err error) {
	slog.Debug("add class", cmd)
	return
}

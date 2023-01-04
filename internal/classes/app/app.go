package app

import (
	"github.com/jbaikge/boneless/internal/classes/app/command"
)

type Commands struct {
	AddClass command.AddClassHandler
	// UpdateClass command.UpdateClassHandler
}

type Queries struct {
	// AllClasses query.AllClassesHandler
	// SingleClass query.SingleClassHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}

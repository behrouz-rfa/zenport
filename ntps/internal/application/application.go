package application

import (
	"context"
	"zenport/ntps/internal/application/commands"
	"zenport/ntps/internal/domain"
)

type (
	App interface {
		Commands
	}
	Commands interface {
		CreateTime(ctx context.Context, cmd commands.CreateTime) (*domain.Time, error)
	}

	Application struct {
		appCommands
	}

	appCommands struct {
		commands.CreateTimeHandler
	}
)

var _ App = (*Application)(nil)

func New(times domain.TimeRepository) *Application {
	return &Application{
		appCommands: appCommands{
			CreateTimeHandler: commands.NewCreateTimeHandler(times),
		},
	}

}

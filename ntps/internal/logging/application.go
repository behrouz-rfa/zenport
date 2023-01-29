package logging

import (
	"context"
	"github.com/rs/zerolog"
	"zenport/ntps/internal/application"
	"zenport/ntps/internal/application/commands"
	"zenport/ntps/internal/domain"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

func LogApplicationAccess(application application.App, logger zerolog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) GetTime(ctx context.Context, request commands.CreateTime) (t *domain.Time, err error) {
	a.logger.Info().Msg("--> Times.GetTime")
	defer func() { a.logger.Info().Err(err).Msg("<-- Times.GetTime") }()
	return a.App.CreateTime(ctx, request)
}

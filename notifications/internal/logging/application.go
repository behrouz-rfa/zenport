package logging

import (
	"context"

	"github.com/rs/zerolog"

	"zenport/notifications/internal/application"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

var _ application.App = (*Application)(nil)

func LogApplicationAccess(application application.App, logger zerolog.Logger) Application {
	return Application{
		App:    application,
		logger: logger,
	}
}

func (a Application) NotifyTimeCreated(ctx context.Context, notify application.TimeCreated) (err error) {
	a.logger.Info().Msg("--> Notifications.NotifyTimeCreated")
	defer func() { a.logger.Info().Err(err).Msg("<-- Notifications.NotifyTimeCreated") }()
	return a.App.NotifyTimeCreated(ctx, notify)
}

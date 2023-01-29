package logging

import (
	"context"
	"github.com/rs/zerolog"
	"zenport/gates/internal/application"
	"zenport/gates/internal/domain"
)

type Application struct {
	application.App
	logger zerolog.Logger
}

func LogApplicationAccess(app application.App, logger zerolog.Logger) Application {
	return Application{App: app, logger: logger}
}

func (a Application) GetTime(ctx context.Context, request application.TimeRequest) (t *domain.Time, err error) {
	a.logger.Info().Msg("--> Gates.GetTime")
	defer func() { a.logger.Info().Err(err).Msg("<-- Gates.GetTime") }()
	return a.App.GetTime(ctx, request)
}

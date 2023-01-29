package application

import (
	"context"
)

type (
	TimeCreated struct {
		Id   string
		Time string
	}

	App interface {
		NotifyTimeCreated(ctx context.Context, notify TimeCreated) error
	}

	Application struct {
		ntps NtpCacheRepository
	}
)

func (a Application) NotifyTimeCreated(ctx context.Context, notify TimeCreated) error {

	return nil
}

var _ App = (*Application)(nil)

func New(nptCachRepo NtpCacheRepository) *Application {
	return &Application{
		ntps: nptCachRepo,
	}
}

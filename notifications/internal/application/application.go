package application

import (
	"context"
)

type (
	OrderCreated struct {
		OrderID    string
		CustomerID string
	}

	OrderCanceled struct {
		OrderID    string
		CustomerID string
	}

	OrderReady struct {
		OrderID    string
		CustomerID string
	}

	App interface {
		NotifyTimeCreated(ctx context.Context, notify OrderCreated) error
	}

	Application struct {
		customers NtpCacheRepository
	}
)

func (a Application) NotifyTimeCreated(ctx context.Context, notify OrderCreated) error {

	return nil
}

var _ App = (*Application)(nil)

func New(customers NtpCacheRepository) *Application {
	return &Application{
		customers: customers,
	}
}

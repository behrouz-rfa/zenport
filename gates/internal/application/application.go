package application

import (
	"context"
	"zenport/gates/internal/domain"
)

type (
	TimeRequest struct {
		Ask string
	}

	App interface {
		GetTime(ctx context.Context, r TimeRequest) (*domain.Time, error)
	}

	Application struct {
		customers domain.TimeRepository
		ntp       domain.NtpRepository
	}
)

var _ App = (*Application)(nil)

func NewApplication(customers domain.TimeRepository, ntps domain.NtpRepository) *Application {
	return &Application{customers: customers, ntp: ntps}
}

func (a Application) GetTime(ctx context.Context, r TimeRequest) (*domain.Time, error) {

	t, err := a.ntp.FetchTime(ctx, r.Ask)
	if err != nil {
		return nil, err
	}

	return domain.NewTime(t), nil
}

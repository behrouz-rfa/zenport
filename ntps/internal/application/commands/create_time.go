package commands

import (
	"context"
	"zenport/ntps/internal/domain"
)

type (
	CreateTime struct {
		ID   string
		Time string
	}

	CreateTimeHandler struct {
		stores domain.TimeRepository
	}
)

func NewCreateTimeHandler(stores domain.TimeRepository) CreateTimeHandler {
	return CreateTimeHandler{
		stores: stores,
	}
}

func (h CreateTimeHandler) CreateTime(ctx context.Context, cmd CreateTime) (*domain.Time, error) {
	store, err := domain.CreateTime(cmd.ID, cmd.Time)
	if err != nil {
		return nil, err
	}

	err = h.stores.Save(ctx, store)
	if err != nil {
		return nil, err
	}
	return store, nil
}

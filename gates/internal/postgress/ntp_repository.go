package postgres

import (
	"context"
	"database/sql"
	"zenport/gates/internal/domain"
)

type TimeRepository struct {
	tableName string
	db        *sql.DB
}

func (t TimeRepository) Load(ctx context.Context, user string) (*domain.Time, error) {
	//TODO implement me
	panic("implement me")
}

func NewTimeRepository(tableName string, db *sql.DB) *TimeRepository {
	return &TimeRepository{tableName: tableName, db: db}
}

var _ domain.TimeRepository = (*TimeRepository)(nil)

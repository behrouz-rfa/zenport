package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"zenport/ntps/internal/domain"
)

type NtpRepository struct {
	tableName string
	db        *sql.DB
}

func (r NtpRepository) AddTime(ctx context.Context, id, time string) error {
	const query = "INSERT INTO %s (id, time) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, r.table(query), id, time)

	return err
}

func NewNtpRepository(tableName string, db *sql.DB) NtpRepository {
	return NtpRepository{tableName: tableName, db: db}
}

var _ domain.NtpRepository = (*NtpRepository)(nil)

func (r NtpRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

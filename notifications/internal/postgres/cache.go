package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"zenport/notifications/internal/application"
)

type NtpCacheRepository struct {
	tableName string
	db        *sql.DB
}

func (n NtpCacheRepository) ShowRequest(ctx context.Context, id, time string) error {
	fmt.Println("Event Added to cache", id, time)
	const query = "INSERT INTO %s (id ,time) VALUES($1,$1)"
	_, err := n.db.ExecContext(ctx, n.table(query), id, time)
	if err != nil {
		return err
	}
	return nil
}

var _ application.NtpCacheRepository = (*NtpCacheRepository)(nil)

func NewNtpCacheRepository(tableName string, db *sql.DB) NtpCacheRepository {
	return NtpCacheRepository{tableName: tableName, db: db}
}

func (n NtpCacheRepository) table(query string) string {
	return fmt.Sprintf(query, n.tableName)
}

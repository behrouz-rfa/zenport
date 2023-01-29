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
	fmt.Println("Event Added to cache")
	return nil
}

var _ application.NtpCacheRepository = (*NtpCacheRepository)(nil)

func NewNtpCacheRepository(tableName string, db *sql.DB) NtpCacheRepository {
	return NtpCacheRepository{tableName: tableName, db: db}
}

package application

import "context"

type NtpCacheRepository interface {
	ShowRequest(ctx context.Context, id, time string) error
}

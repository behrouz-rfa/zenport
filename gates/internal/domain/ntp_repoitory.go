package domain

import "context"

type NtpRepository interface {
	FetchTime(ctx context.Context, request string) (string, error)
}

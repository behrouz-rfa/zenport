package domain

import (
	"context"
)

type NtpStore struct {
	ID   string
	Time string
}

type NtpRepository interface {
	AddTime(ctx context.Context, id, time string) error
}

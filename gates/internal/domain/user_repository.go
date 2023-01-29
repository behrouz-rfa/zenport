package domain

import "context"

type TimeRepository interface {
	Load(ctx context.Context, user string) (*Time, error)
}

package domain

import "context"

type TimeRepository interface {
	Save(ctx context.Context, t *Time) error
}
